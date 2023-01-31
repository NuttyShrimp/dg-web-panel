/// <reference types="vite/client" />

interface ImportMetaEnv {
  VITE_BACKEND_OVERWRITE: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
