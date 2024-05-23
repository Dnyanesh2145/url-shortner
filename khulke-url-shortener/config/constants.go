package config

const (
	CONFIG_LOGFORMAT = "${time} <----> ${ip} <----> ${locals:reqid} <----> ${locals:userId} <----> ${method}" +
		" <----> ${path} <----> ${status} <----> ${bytesSent} <----> ${latency} <----> ${ua} <----> ${locals:device}" +
		" <----> ${locals:latlong} <----> ${referrer} <----> ${locals:msg}"
	CONFIG_LOGTIME_FORMAT = "2006-01-02T15:04:05.000"
	CONFIG_DEFAULT_REQID  = "000000"
	CONFIG_DEFAULT_UID    = "ffffff"
	DEV_ENVIORMENT        = "dev"
	SHORT_URL_KEY         = "shortened_url_keys"
)
