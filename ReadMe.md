# Go Proxy Dialers

This repo contains Dialers http connect and https for golang.org/x/net/proxy

# Using HTTP Connect Dialer

import the connect package first and follow the below examples

## With HTTP Client
> **Note:** for http requests (its preferable to HTTP_PROXY and HTTPS_PROXY env vars in this case)
```golang
	uri, err := url.Parse("http://proxy_url:proxy_port")
	if err != nil {
		panic(err)
	}

	dialer, err := proxy.FromURL(uri, proxy.Direct)
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy:       nil,
			DialContext: dialer.Dial,
		},
	}

```

## With GRPC 

In this case we need to make use of the **grpc.WithContextDialer** DialOption

```golang
grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
    uri, err := url.Parse(httpproxy.FromEnvironment().HTTPProxy)
    if err != nil {
        return nil, err
    }

    dialer, err := proxy.FromURL(uri, proxy.Direct)
    if err != nil {
        return nil, err
    }
    return dialer.Dial("tcp", addr)
}),
```

# Using HTTPS Dialer

import the https package and follow the below example

```golang
uri, err := url.Parse(httpproxy.FromEnvironment().HTTPProxy)
    if err != nil {
        return nil, err
    }

    dialer, err := proxy.FromURL(uri, https.HTTPSDialer)
    if err != nil {
        return nil, err
    }
```