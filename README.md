
This repository contains a list of of HTTP user-agents used by robots, crawlers, and spiders. I regularly maintain this list based on my own logs. 

If you are using Ruby, [Voight-Kampff](https://github.com/biola/Voight-Kampff) and [isbot](https://github.com/Hentioe/isbot) provide  libraries for accessing this data.

Other systems for spotting robots, crawlers, and spiders that you may want to consider include [isBot](https://github.com/gorangajic/isbot) (Node.JS), [Crawler-Detect](https://github.com/JayBizzle/Crawler-Detect) (PHP), [BrowserDetector](https://github.com/mimmi20/BrowserDetector) (PHP), and [browscap](https://github.com/browscap/browscap) (JSON files).

## License

The list is under a [MIT License](https://opensource.org/licenses/MIT). The versions prior to Nov 7, 2016 were under a [CC-SA](http://creativecommons.org/licenses/by-sa/3.0/) license.

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
      "url": "http://moz.com/help/pro/what-is-rogerbot-"
    }


--Martin
