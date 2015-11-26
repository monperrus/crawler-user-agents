
The Github repository `crawler-user-agents` contains a list of of HTTP user-agents used by robots/crawlers/spiders. I regularly maintain this list based on my own logs. I do welcome additions contributed as pull requests. 

The pull requests should:

* contain few additions (say less than 5)
* specify a discriminant relevant syntactic fragment (for example "totobot" and not "Mozilla/5 totobot v20131212.alpha1") 
* contain the pattern (generic regular expression), the discovery date (year/month/day) and the official url of the robot
* result in a valid JSON file (don't forget the comma between items)

Example:

    {
      "pattern": "rogerbot", 
      "addition_date": "2014/02/28", 
      "url": "http://moz.com/help/pro/what-is-rogerbot-"
    }


The list is under a [CC-SA](http://creativecommons.org/licenses/by-sa/3.0/) license.

--Martin

## Create a single RegExp rule for matching crawlers
1. Open your browser console
2. Run

    var botPattern, botRegexp, xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(evt){var xhr = evt.target; if (xhr.readyState==4 && xhr.status==200) {
        var bots = JSON.parse(xhr.responseText);
        botPattern="(";bots.forEach(function(item){botPattern+=item.pattern+'|';});
        botPattern=botPattern.substring(0,botPattern.length-1)+')';
        console.info("The botPattern is %o", botPattern);
        console.log('You can type "botPattern" to show crawler-matching pattern again');
        botRegexp = new RegExp(botPattern, 'i'); // match case insensitive
        console.log('You can use botRegexp to test User-Agent string. Like this:\n\nbotRegexp.test(\'Googlebot/2.1 (+http://www.googlebot.com/bot.html)\')');
    }};
    xhr.open('GET', 'https://raw.githubusercontent.com/monperrus/crawler-user-agents/master/crawler-user-agents.json', true);
    xhr.send();
3. Use botPattern to create new RegExp(botPattern, 'i') instance to .test() against the User-Agent string;

