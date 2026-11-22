declare const crawlerUserAgents: {
	pattern: string
	addition_date?: string
	url?: string
	instances: string[]
	verification?: {
		type: "web-bot-auth"
		directory: string
	}[]
	tags?: string[]
}[]

export default crawlerUserAgents;
