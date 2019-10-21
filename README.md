# shinobi-alerter
a way to alert to shinobi ip camera software


# Config example
`config.json`
```
{
	"server": "https://shinobi.video",
	"apikey": "oEhEEeUWaTLA92xiY9oXq9BGK5QdvE",
	"cameras": [
		{
			"name": "Frontdoor-01",
			"ip": "10.1.1.2",
			"group": "bPBPi7PLzg",
			"region": "front"
		},{
			"name": "Driveway-01",
			"ip": "10.1.1.3",
			"group": "bPBPi7PLzg",
			"region": "g1iU3"
		}
	]
}
```