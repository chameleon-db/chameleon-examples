// Configuración de la API
const API_BASE_URL = 'http://localhost:8080';
let currentUser = null;
let todos = [];

// Elementos del DOM
const loginCard = document.getElementById('loginCard');
const dashboard = document.getElementById('dashboard');
const loginForm = document.getElementById('loginForm');
const logoutBtn = document.getElementById('logoutBtn');
const todoForm = document.getElementById('todoForm');
const todosContainer = document.getElementById('todosContainer');
const loadingSpinner = document.getElementById('loadingSpinner');
const toast = document.getElementById('toast');
const userEmailDisplay = document.getElementById('userEmailDisplay');
const totalTasksEl = document.getElementById('totalTasks');
const completedTasksEl = document.getElementById('completedTasks');
const pendingTasksEl = document.getElementById('pendingTasks');

// Helper: Mostrar Toast
function showToast(message, isError = false) {
    toast.textContent = message;
    toast.style.borderLeftColor = isError ? '#ef4444' : '#f97316';
    toast.classList.add('show');
    setTimeout(() => toast.classList.remove('show'), 3000);
}

function normalizeUuid(value) {
    if (typeof value === 'string') {
        return value;
    }

    if (Array.isArray(value) && value.length === 16 && value.every(n => Number.isInteger(n) && n >= 0 && n <= 255)) {
        const hex = value.map(n => n.toString(16).padStart(2, '0')).join('');
        return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
    }

    return value;
}

function userTodosBasePath() {
    if (!currentUser?.id) {
        throw new Error('Usuario no autenticado');
    }

    return `/users/${currentUser.id}/todos`;
}

// API Calls con manejo de errores
async function apiRequest(endpoint, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };

    const options = {
        method,
        headers,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
        
        // Handle empty responses (204 No Content, etc.)
        const contentType = response.headers.get('content-type');
        let data;
        
        if (contentType && contentType.includes('application/json')) {
            data = await response.json();
        } else {
            data = await response.text();
        }
        
        if (!response.ok) {
            const errorMsg = data?.error || data?.message || data || 'Error en la solicitud';
            throw new Error(errorMsg);
        }

        // Handle wrapped responses
        return data && data.data !== undefined ? data.data : data;
    } catch (error) {
        console.error('API Error:', error);
        showToast(error.message, true);
        throw error;
    }
}

// LOGIN
async function handleLogin(email, password) {
    try {
        const loginData = {
            email: email,
            password: password
        };
        
        console.log('📤 Enviando login:', { email: loginData.email });
        
        const user = await apiRequest('/login', 'POST', loginData);
        const userId = normalizeUuid(user?.id);
        
        if (!user || !userId || typeof userId !== 'string') {
            throw new Error('Respuesta inválida del servidor');
        }
        
        currentUser = {
            id: userId,
            email: user.email,
        };
        
        localStorage.setItem('user', JSON.stringify(currentUser));
        showToast(`¡Bienvenido, ${currentUser.email}!`);
        return currentUser;
        
    } catch (error) {
        console.error('Login error:', error);
        showToast(error.message || 'Error de autenticación', true);
        throw error;
    }
}

// Cargar TODOs del usuario
async function loadTodos() {
    if (!currentUser) return;
    
    try {
        loadingSpinner.style.display = 'block';
        todosContainer.innerHTML = '';
        
        const data = await apiRequest(`${userTodosBasePath()}/`);

        todos = Array.isArray(data)
            ? data.map(todo => ({
                ...todo,
                id: normalizeUuid(todo.id),
                user_id: normalizeUuid(todo.user_id),
            }))
            : [];
        renderTodos();
        updateStats();
        
    } catch (error) {
        console.error('Load todos error:', error);
        todosContainer.innerHTML = '<p style="color: #a0a0a0; text-align: center;">Error al cargar tareas</p>';
    } finally {
        loadingSpinner.style.display = 'none';
    }
}

// Renderizar lista de TODOs
function renderTodos() {
    if (todos.length === 0) {
        todosContainer.innerHTML = '<p style="color: #a0a0a0; text-align: center; padding: 40px;">No hay tareas aún. ¡Añade una!</p>';
        return;
    }

    todosContainer.innerHTML = todos.map(todo => `
        <div class="todo-item" data-id="${todo.id}">
            <div class="todo-content">
                <div class="todo-check ${todo.completed ? 'completed' : ''}" onclick="toggleTodo('${todo.id}')">
                    ${todo.completed ? '<i class="fas fa-check"></i>' : ''}
                </div>
                <span class="todo-title ${todo.completed ? 'completed' : ''}">${todo.title || 'Sin título'}</span>
            </div>
            <div class="todo-actions">
                <button onclick="deleteTodo('${todo.id}')"><i class="fas fa-trash"></i></button>
            </div>
        </div>
    `).join('');
}

// Actualizar estadísticas
function updateStats() {
    const total = todos.length;
    const completed = todos.filter(t => t.completed).length;
    const pending = total - completed;

    totalTasksEl.textContent = total;
    completedTasksEl.textContent = completed;
    pendingTasksEl.textContent = pending;
}

async function createTodo(title) {
    if (!currentUser) return;
    
    try {
        console.log('📤 Creating todo:', { title, userId: currentUser.id });
        
        const newTodo = await apiRequest(`${userTodosBasePath()}/`, 'POST', {
            title: title,
            description: "Default description",
        });

        console.log('✅ Todo created:', newTodo);

        if (newTodo) {
            newTodo.id = normalizeUuid(newTodo.id);
            newTodo.user_id = normalizeUuid(newTodo.user_id);
            newTodo.completed = newTodo.completed || false;
        }
        
        todos.unshift(newTodo);
        renderTodos();
        updateStats();
        showToast('Tarea creada');
        document.getElementById('todoTitle').value = '';
    } catch (error) {
        console.error('❌ Create todo error:', error);
    }
}

window.toggleTodo = async function(todoId) {
    const todo = todos.find(t => t.id == todoId);
    if (!todo) return;
    
    try {
        console.log('📤 Toggling todo:', todoId);
        
        await apiRequest(`${userTodosBasePath()}/${todoId}/toggle`, 'PATCH');
        
        // Actualizar estado local
        todo.completed = !todo.completed;
        renderTodos();
        updateStats();
        showToast(`Tarea ${todo.completed ? 'completada' : 'pendiente'}`);
    } catch (error) {
        console.error('❌ Toggle error:', error);
    }
}

window.deleteTodo = async function(todoId) {
    if (!confirm('¿Eliminar tarea?')) return;
    
    try {
        console.log('📤 Deleting todo:', todoId);
        
        await apiRequest(`${userTodosBasePath()}/${todoId}`, 'DELETE');
        
        todos = todos.filter(t => t.id != todoId);
        renderTodos();
        updateStats();
        showToast('Tarea eliminada');
    } catch (error) {
        console.error('❌ Delete error:', error);
    }
}

// Logout
function logout() {
    currentUser = null;
    localStorage.removeItem('user');
    loginCard.style.display = 'block';
    dashboard.style.display = 'none';
    todos = [];
    showToast('Sesión cerrada');
}

// Inicializar desde localStorage
function initFromStorage() {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
        try {
            currentUser = JSON.parse(storedUser);
            const normalizedId = normalizeUuid(currentUser?.id);

            if (currentUser && normalizedId && typeof normalizedId === 'string') {
                currentUser.id = normalizedId;
                localStorage.setItem('user', JSON.stringify(currentUser));
                loginCard.style.display = 'none';
                dashboard.style.display = 'block';
                userEmailDisplay.textContent = currentUser.email;
                loadTodos();
            } else {
                localStorage.removeItem('user');
            }
        } catch (e) {
            localStorage.removeItem('user');
        }
    }
}

// Event Listeners
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    const loginBtn = document.getElementById('loginBtn');
    loginBtn.disabled = true;
    loginBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Entrando...';
    
    try {
        await handleLogin(email, password);
        loginCard.style.display = 'none';
        dashboard.style.display = 'block';
        userEmailDisplay.textContent = currentUser.email;
        loadTodos();
    } catch (error) {
        // Error ya manejado en handleLogin
    } finally {
        loginBtn.disabled = false;
        loginBtn.innerHTML = '<span>Login</span><i class="fas fa-arrow-right"></i>';
    }
});

logoutBtn.addEventListener('click', logout);

todoForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const input = document.getElementById('todoTitle');
    const title = input.value.trim();
    if (title) {
        createTodo(title);
    }
});

// Inicializar
initFromStorage();