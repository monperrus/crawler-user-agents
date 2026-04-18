// Example:
// {
// 	"pattern": "rogerbot",
// 	"addition_date": "2014/02/28",
// 	"url": "http://moz.com/help/pro/what-is-rogerbot-",
// 	"instances" : ["rogerbot/2.3 example UA"],
// 	"tags": ["seo"]
// }

declare const crawlerUserAgents: {
	pattern: string
	addition_date?: string
	url?: string
	instances: string[]
	tags?: string[]
}[]

export = crawlerUserAgents;
