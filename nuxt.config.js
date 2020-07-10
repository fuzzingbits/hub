export default {
	mode: "spa",
	dir: {
		layouts: "ui/layouts",
		middleware: "ui/middleware",
		pages: "ui/pages",
		static: "ui/static",
		store: "ui/store",
	},
	buildModules: ["@nuxt/typescript-build"],
};
