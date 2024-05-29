package specs

import (
	"errors"
	netur "net/url"
	"strconv"
	"strings"
)

var invalidFormatError = errors.New("url: invalid format")

func ParseUrl(url string) (*Url, error) {
	obj := &Url{}
	if len(url) > 0 {
		var i, mark, step int
		end := len(url) - 1
		for i < end {
			switch url[i] {
			case '/':
				if step != 5 {
					if i != 0 {
						switch step { // read as 'path'
							default:
								return nil, invalidFormatError
	
							case 0, 3: // from 'host'
								if i - mark < 1 {
									return nil, invalidFormatError
								}
								obj.Host = url[mark:i]
			
							case 4: // from 'port'
								err := obj.setPort(url[mark:i])
								if err != nil {
									return nil, errors.New("url -> port: " + err.Error())
								}
						}
					}

					step = 5 // goto 'path'
					mark = i
				}
			case ':':
				if i < 1 { // ensure host:port format
					return nil, invalidFormatError
				} else if i+2 < end && 
					url[i+1] == '/' && url[i+2] == '/' { // read as 'scheme'
					
					if step != 0 && i < 1 {
						return nil, invalidFormatError
					}
					step = 3 // goto 'host'
					obj.Scheme = url[:i]
					i += 3; mark = i
					continue
				} else if step == 0 || step == 3 { // read as 'host'
					if i - mark < 1 {
						return nil, invalidFormatError
					}
					obj.Host = url[mark:i]
					step = 4 // goto 'port'
					i++; mark = i
					continue
				}
			case '@':
				if step == 4 { // from 'port', read as 'password'
					if i - mark < 1 {
						return nil, invalidFormatError
					}
					obj.Username = obj.Host
					obj.Password = url[mark:i]
					obj.Host = ""
					step = 3 // goto 'host'
					i++; mark = i
					continue
				} else {
					return nil, invalidFormatError
				}
			case '?':
				if step == 5 { // from 'path', read as 'query'
					if i - mark < 1 {
						return nil, invalidFormatError
					}
					obj.Path = url[mark:i]
					step = 6 // goto 'query'
					i++; mark = i
					continue
				} else {
					return nil, invalidFormatError
				}
			case '#':
				switch step { // read as 'hash'
					default:
						return nil, invalidFormatError

					case 5: // from 'path'
						obj.Path = url[mark:i]

					case 6: // from 'query'
						obj.query = url[mark:i]
				}
				step = 7 // goto 'hash'
				i++; mark = i
				continue
			}
			i++
		}
		if end - mark < 0 {
			return nil, invalidFormatError
		}
		switch step {
			case 0, 3: // host
				obj.Host = url[mark:]
			
			case 4: // port
				err := obj.setPort(url[mark:])
				if err != nil {
					return nil, errors.New("url -> port: " + err.Error())
				}

			case 5: // path
				obj.Path = url[mark:]

			case 6: // query
				obj.query = url[mark:]

			case 7: // hash
				obj.Hash = url[mark:]

			default:
				return nil, invalidFormatError
		}
	}
	if len(obj.Path) > 2 {
		if path, err := netur.PathUnescape(obj.Path); err == nil {
			obj.Path = path
		}
	}
	return obj, nil
}

type Url struct {
	Scheme, Username, Password, 
	Host, Path, query, Hash string
	Port uint16

	queryParams Query
}

func (url *Url) setPort(val string) error {
	num, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return err
	}
	url.Port = uint16(num)
	return nil
}

func (url *Url) Query() string {
	return url.query
}

func (url *Url) SetQuery(query string) {
	url.queryParams = nil
	url.query = query
}

func (url *Url) QueryParams() Query {
	if url.queryParams != nil {
		return url.queryParams
	}
	
	query, err := ParseQuery(url.query)
	if err != nil {
		query = Query{}
	}
	url.queryParams = query
	return query
}

func (url *Url) String() string {
	var builder strings.Builder

	if len(url.Host) > 0 {
		if len(url.Scheme) > 0 {
			builder.WriteString(url.Scheme)
			builder.WriteString("://")
		}
		if len(url.Username) > 0 && len(url.Password) > 0 {
			builder.WriteString(url.Username)
			builder.WriteByte(':')
			builder.WriteString(url.Password)
			builder.WriteByte('@')
		}

		builder.WriteString(url.Host)

		if url.Port > 0 {
			builder.WriteByte(':')
			builder.Write(strconv.AppendUint(nil, uint64(url.Port), 10))
		}
	}

	if len(url.Path) > 0 {
		builder.WriteString(url.Path)

		if len(url.query) > 0 {
			builder.WriteByte('?')
			builder.WriteString(url.query)
		}
		if len(url.Hash) > 0 {
			builder.WriteByte('#')
			builder.WriteString(url.Hash)
		}
	}

	return builder.String()
}
