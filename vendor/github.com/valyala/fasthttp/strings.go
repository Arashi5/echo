package fasthttp

var (
	defaultServerName  = "fasthttp"
	defaultUserAgent   = "fasthttp"
	defaultContentType = []byte("text/plain; charset=utf-8")
)

var (
	strSlash                    = []byte("/")
	strSlashSlash               = []byte("//")
	strSlashDotDot              = []byte("/..")
	strSlashDotSlash            = []byte("/./")
	strSlashDotDotSlash         = []byte("/../")
	strBackSlashDotDot          = []byte(`\..`)
	strBackSlashDotBackSlash    = []byte(`\.\`)
	strSlashDotDotBackSlash     = []byte(`/..\`)
	strBackSlashDotDotBackSlash = []byte(`\..\`)
	strCRLF                     = []byte("\r\n")
	strHTTP                     = []byte("http")
	strHTTPS                    = []byte("https")
	strHTTP11                   = []byte("HTTP/1.1")
	strColon                    = []byte(":")
	strColonSlashSlash          = []byte("://")
	strColonSpace               = []byte(": ")
	strCommaSpace               = []byte(", ")
	strGMT                      = []byte("GMT")

	strResponseContinue = []byte("HTTP/1.1 100 Continue\r\n\r\n")

	strExpect             = []byte(HeaderExpect)
	strConnection         = []byte(HeaderConnection)
	strContentLength      = []byte(HeaderContentLength)
	strContentType        = []byte(HeaderContentType)
	strDate               = []byte(HeaderDate)
	strHost               = []byte(HeaderHost)
	strReferer            = []byte(HeaderReferer)
	strServer             = []byte(HeaderServer)
	strTransferEncoding   = []byte(HeaderTransferEncoding)
	strContentEncoding    = []byte(HeaderContentEncoding)
	strAcceptEncoding     = []byte(HeaderAcceptEncoding)
	strUserAgent          = []byte(HeaderUserAgent)
	strCookie             = []byte(HeaderCookie)
	strSetCookie          = []byte(HeaderSetCookie)
	strLocation           = []byte(HeaderLocation)
	strIfModifiedSince    = []byte(HeaderIfModifiedSince)
	strLastModified       = []byte(HeaderLastModified)
	strAcceptRanges       = []byte(HeaderAcceptRanges)
	strRange              = []byte(HeaderRange)
	strContentRange       = []byte(HeaderContentRange)
	strAuthorization      = []byte(HeaderAuthorization)
	strTE                 = []byte(HeaderTE)
	strTrailer            = []byte(HeaderTrailer)
	strMaxForwards        = []byte(HeaderMaxForwards)
	strProxyConnection    = []byte(HeaderProxyConnection)
	strProxyAuthenticate  = []byte(HeaderProxyAuthenticate)
	strProxyAuthorization = []byte(HeaderProxyAuthorization)
	strWWWAuthenticate    = []byte(HeaderWWWAuthenticate)
	strVary               = []byte(HeaderVary)

	strCookieExpires        = []byte("expires")
	strCookieDomain         = []byte("domain")
	strCookiePath           = []byte("path")
	strCookieHTTPOnly       = []byte("HttpOnly")
	strCookieSecure         = []byte("secure")
	strCookieMaxAge         = []byte("max-age")
	strCookieSameSite       = []byte("SameSite")
	strCookieSameSiteLax    = []byte("Lax")
	strCookieSameSiteStrict = []byte("Strict")
	strCookieSameSiteNone   = []byte("None")

	strClose               = []byte("close")
	strGzip                = []byte("gzip")
	strBr                  = []byte("br")
	strDeflate             = []byte("deflate")
	strKeepAlive           = []byte("keep-alive")
	strUpgrade             = []byte("Upgrade")
	strChunked             = []byte("chunked")
	strIdentity            = []byte("identity")
	str100Continue         = []byte("100-continue")
	strPostArgsContentType = []byte("application/x-www-form-urlencoded")
	strDefaultContentType  = []byte("application/octet-stream")
	strMultipartFormData   = []byte("multipart/form-data")
	strBoundary            = []byte("boundary")
	strBytes               = []byte("bytes")
	strBasicSpace          = []byte("Basic ")

	strApplicationSlash = []byte("application/")
	strImageSVG         = []byte("image/svg")
	strImageIcon        = []byte("image/x-icon")
	strFontSlash        = []byte("font/")
	strMultipartSlash   = []byte("multipart/")
	strTextSlash        = []byte("text/")
)