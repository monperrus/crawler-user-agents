# crawler-user-agents

This repository contains a list of of HTTP user-agents used by robots, crawlers, and spiders as in single JSON file.

## Install

### Direct download

Download the [`crawler-user-agents.json` file](https://raw.githubusercontent.com/monperrus/crawler-user-agents/master/crawler-user-agents.json) from this repository directly.

### Npm / Yarn

crawler-user-agents is deployed on npmjs.com: <https://www.npmjs.com/package/crawler-user-agents>

To use it using npm or yarn:

```sh
npm install --save crawler-user-agents
# OR
yarn add crawler-user-agents
```

In Node.js, you can `require` the package to get an array of crawler user agents.

```js
const crawlers = require('crawler-user-agents');
console.log(crawlers);
```

## Usage

Each `pattern` is a regular expression. It should work out-of-the-box wih your favorite regex library:

* JavaScript: `if (RegExp(entry.pattern).test(req.headers['user-agent']) { ... }`
* PHP: add a slash before and after the pattern: `if (preg_match('/'.$entry['pattern'].'/', $_SERVER['HTTP_USER_AGENT'])): ...`
* Python: `if re.search(entry['pattern'], ua): ...`

## Contributing

I do welcome additions contributed as pull requests.

The pull requests should:

* contain a single addition
* specify a discriminant relevant syntactic fragment (for example "totobot" and not "Mozilla/5 totobot v20131212.alpha1")
* contain the pattern (generic regular expression), the discovery date (year/month/day) and the official url of the robot
* result in a valid JSON file (don't forget the comma between items)

Example:

    {
      "pattern": "rogerbot",
      "addition_date": "2014/02/28",
      "url": "http://moz.com/help/pro/what-is-rogerbot-",
      "instances" : ["rogerbot/2.3 example UA"]
    }

## License

The list is under a [MIT License](https://opensource.org/licenses/MIT). The versions prior to Nov 7, 2016 were under a [CC-SA](http://creativecommons.org/licenses/by-sa/3.0/) license.

## Related work

There are a few wrapper libraries that use this data to detect bots:

 * [Voight-Kampff](https://github.com/biola/Voight-Kampff) (Ruby)
 * [isbot](https://github.com/Hentioe/isbot) (Ruby)
 * [crawlers](https://github.com/Olical/crawlers) (Clojure)
 * [crawlerflagger](https://godoc.org/go.kelfa.io/kelfa/pkg/crawlerflagger) (Go)
 * [isBot](https://github.com/omrilotan/isbot) (Node.JS)

Other systems for spotting robots, crawlers, and spiders that you may want to consider are:

 * [Crawler-Detect](https://github.com/JayBizzle/Crawler-Detect) (PHP)
 * [BrowserDetector](https://github.com/mimmi20/BrowserDetector) (PHP)
 * [browscap](https://github.com/browscap/browscap) (JSON files)
