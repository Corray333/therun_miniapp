declare module '*.vue' {
  import Vue from 'vue'
  export default Vue
}

interface ImportMeta {
  globEager: (pattern: string) => Record<string, () => Promise<any>>;
}