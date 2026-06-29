# Agent Guidelines

Conventions for AI agents working in this repository.

## TypeScript

### Prefer `type` over `interface`

Use `type` aliases for object shapes, props, and API payloads. Do not introduce new `interface` declarations.

```typescript
// ❌ BAD
interface UploadFormProps {
  isLoading: boolean;
  handleSubmit: () => void;
}

// ✅ GOOD
type UploadFormProps = {
  isLoading: boolean;
  handleSubmit: () => void;
};
```

Apply this in all `*.ts` and `*.tsx` files under `frontend/src/`.

### Use semantic variable names

In `frontend/src/`, avoid single-character variable names (for example `k`, `i`, `n`). Names should describe intent or role, not just type or position.

```typescript
// ❌ BAD
const k = 1024;
const i = Math.floor(Math.log(bytes) / Math.log(k));

// ✅ GOOD
const bytesPerUnit = 1024;
const unitIndex = Math.floor(Math.log(bytes) / Math.log(bytesPerUnit));
```

Short names are fine in very small scopes when the meaning is obvious from context (for example `index` in `items.map((item, index) => ...)`), but prefer descriptive names for anything that carries business or formatting logic.

### Prefer absolute imports with the `@/` alias

In `frontend/src/`, avoid relative imports (`../`, `./`). Use the `@/` webpack/TypeScript alias, which resolves to `src/`.

```typescript
// ❌ BAD
import logoLockupLight from "../assets/logo-lockup-light.svg";
import { useTheme } from "../context/ThemeContext";

// ✅ GOOD
import logoLockupLight from "@/assets/logo-lockup-light.svg";
import { useTheme } from "@/context/ThemeContext";
```

The alias is configured in `frontend/vite.config.ts` and `frontend/tsconfig.json`.
