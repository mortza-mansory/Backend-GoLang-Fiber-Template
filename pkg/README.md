# pkg

The `pkg/` directory is reserved for code intentionally designed to be
imported by external applications.

**Keep this minimal.** Do not move internal application code here. Internal
packages live under `internal/`, which Go prevents other modules from
importing.

Only place code in `pkg/` if it is a genuinely reusable library that you want
third parties to consume. For most projects, this directory can stay empty.