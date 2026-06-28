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
