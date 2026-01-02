import { defineConfig } from "vitepress";
import llmstxt, { copyOrDownloadAsMarkdownButtons } from "vitepress-plugin-llms";
// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Marvin",
  description: "Documentation QA",
  // Configure favicon
  head: [
    ["link", { rel: "icon", type: "image/svg+xml", href: "/favicon.svg" }],
  ],

  vite: {
    plugins: [llmstxt()],
  },

  markdown: {
    config(md) {
      md.use(copyOrDownloadAsMarkdownButtons);
    },
  },

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Examples", link: "/markdown-examples" },
    ],

    search: {
      provider: 'local'
    },

    sidebar: [
      {
        text: 'Installation',
        link: 'installation'
        },
        {
        text: 'Getting Started',
        link: 'getting-started'
        },
      {
        text: "Reference",
        items: [{ text: "CLI", link: "/reference/cli" }],
      },
      {
        text: "Contributing",
        collapsed: true,
        items: [
          { text: "CLI", link: "/contribute/cli/index" },
          { text: "Web App", link: "/contribute/web/index" },
          { text: "Documentation", link: "/contribute/docs/index" }
        ],
      },
      {
        text: "Examples",
        items: [
          { text: "Markdown Examples", link: "/markdown-examples" },
          { text: "Runtime API Examples", link: "/api-examples" },
        ],
      },
    ],

    socialLinks: [{ icon: "github", link: "https://github.com/svx/marvin" }],
  },
});
