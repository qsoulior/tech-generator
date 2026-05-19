import "vfonts/Roboto.css"

import { createApp } from "vue"
import { createPinia } from "pinia"
import { config, type CodeMirrorExtension } from "md-editor-v3"
import RU from "@vavt/cm-extension/dist/locale/ru"

import App from "./App.vue"
import router from "./router"
import { templatePlaceholderExtension } from "./utils/templatePlaceholder"
import "./assets/template-placeholder.css"

config({
  editorConfig: {
    languageUserDefined: {
      ru: RU,
    },
  },
  codeMirrorExtensions: (extensions: CodeMirrorExtension[]) => [
    ...extensions,
    { type: "template-placeholder", extension: templatePlaceholderExtension },
  ],
})

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount("#app")
