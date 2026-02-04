# ChameleonDB Examples

Real-world examples showing how to use ChameleonDB in different scenarios.

Each example is self-contained with schema, seed data, and documentation.

## Getting Started

### Prerequisites

1. Install ChameleonDB CLI:
```bash
   curl -sSL https://chameleondb.dev/install.sh | sh
```

2. Have PostgreSQL running:
```bash
   # macOS
   brew services start postgresql
   
   # Linux
   sudo systemctl start postgresql
   
   # Docker
   docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16
```

### Run an Example
```bash
git clone https://github.com/chameleon-db/chameleondb-examples
cd chameleondb-examples/01-hello-world
chameleon validate
createdb hello_chameleon
chameleon migrate --apply
psql -U postgres hello_chameleon < seed.sql
```

## Examples

### Beginner

**[01-hello-world](./01-hello-world)** — 5 minutes  
Single entity. Learn the basics: schema, migration, queries.

Concepts: entities, fields, constraints, basic queries

### Intermediate

**[02-blog](./02-blog)** — 15 minutes *(coming soon)*  
User, Post, Comment with relations.

Concepts: HasMany, BelongsTo, eager loading, nested includes

**[03-ecommerce](./03-ecommerce)** — 30 minutes *(coming soon)*  
Product catalog, orders, cart.

Concepts: decimals, complex domains, calculated fields

### Advanced

**[04-saas-multitenant](./04-saas-multitenant)** — 45 minutes *(coming soon)*  
Multi-tenant SaaS with tenant isolation.

Concepts: tenant scoping, row-level security patterns

**[05-social-network](./05-social-network)** — 60 minutes *(coming soon)*  
Followers, posts, likes, comments.

Concepts: many-to-many, graph navigation, complex filters

## Step-by-Step Tutorial

New to ChameleonDB? Follow our interactive tutorial:

**[Tutorial: Building a Task Manager](./tutorial)** *(coming soon)*

Learn by building a real app, step by step.

## Contributing

Have an example to share? PRs welcome!

1. Fork this repo
2. Create your example in a new directory
3. Follow the structure of existing examples
4. Submit a PR

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## Support

- [Documentation](https://github.com/chameleon-db/chameleondb)
- [Discord Community](https://discord.gg/tyZNY2xmr)
- [GitHub Discussions](https://github.com/chameleon-db/chameleondb/discussions)

## License

All examples are MIT licensed. Use them however you want.
