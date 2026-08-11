import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'

const vueRecommended = pluginVue.configs['flat/recommended']

export default [
  {
    ignores: ['dist', 'node_modules', 'public', 'coverage'],
  },
  {
    files: ['**/*.vue'],
    languageOptions: {
      globals: { ...globals.browser },
      parser: pluginVue.parser,
      parserOptions: {
        parser: '@typescript-eslint/parser',
      },
    },
  },
  ...(Array.isArray(vueRecommended) ? vueRecommended : [vueRecommended]),
  {
    rules: {
      'vue/multi-word-component-names': 'off',
      'vue/html-self-closing': 'off',
      'vue/max-attributes-per-line': 'off',
    },
  },
]
