## 概览

首先同步下项目概况：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/7_api_1.png)

上篇文章分享了，路由中间件 - Jaeger 链路追踪（实战篇），文章反响真是出乎意料， 「Go中国」 公众号也转发了，有很多朋友加我好友交流，直呼我大神，其实我哪是什么大神，只不过在本地实践了而已，对于 Go 语言的使用，我还是个新人，在这里感谢大家的厚爱！

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/7_api_2.gif)

这篇文章咱们分享：路由中间件 - 签名验证。

为什么使用签名验证？

这个就不用多说了吧，主要是为了保证接口安全和识别调用方身份，基于这两点，咱们一起设计下签名。

调用方需要申请 App Key 和 App Secret，App Key 用来识别调用方身份，App Secret 用来加密生成签名使用。

当然生成的签名还需要满足以下几点：

- 可变性：每次的签名必须是不一样的。
- 时效性：每次请求的时效性，过期作废。
- 唯一性：每次的签名是唯一的。
- 完整性：能够对传入数据进行验证，防止篡改。

举个例子：

`/api?param_1=xxx&param_2=xxx`，其中 param_1 和 param_2 是两个参数。

如果增加了签名验证，需要再传递几个参数：

- ak 表示App Key，用来识别调用方身份。
- ts 表示时间戳，用来验证接口的时效性。
- sn 表示签名加密串，用来验证数据的完整性，防止数据篡改。

sn 是通过 App Secret 和 传递的参数 进行加密的。

最终传递的参数如下：

`/api?param_1=xxx&param_2=xxx&ak=xxx&ts=xxx&sn=xxx`

在这要说一个调试技巧，ts 和 sn 参数每次都手动生成太麻烦了，当传递 `debug=1` 的时候，会返回 ts 和 sn , 具体看下代码就清楚了。

这篇文章分享三种实现签名的方式，分别是：MD5 组合加密、AES 对称加密、RSA 非对称加密。

废话不多说，进入主题。

## MD5 组合

#### 生成签名

首先，封装一个 Go 的 MD5 方法：

```
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}
```

进行加密：

```
appKey     = "demo"
appSecret  = "xxx"
encryptStr = "param_1=xxx&param_2=xxx&ak="+appKey+"&ts=xxx"

// 自定义验证规则
sn = MD5(appSecret + encryptStr + appSecret)
```

#### 验证签名

通过传递参数，再次生成签名，如果将传递的签名与生成的签名进行对比。

相同，表示签名验证成功。

不同，表示签名验证失败。

#### 中间件 - 代码实现

```
var AppSecret string

// MD5 组合加密
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		utilGin := util.Gin{Ctx: c}

		sign, err := verifySign(c)

		if sign != nil {
			utilGin.Response(-1, "Debug Sign", sign)
			c.Abort()
			return
		}

		if err != nil {
			utilGin.Response(-1, err.Error(), sign)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证签名
func verifySign(c *gin.Context) (map[string]string, error) {
	_ = c.Request.ParseForm()
	req   := c.Request.Form
	debug := strings.Join(c.Request.Form["debug"], "")
	ak    := strings.Join(c.Request.Form["ak"], "")
	sn    := strings.Join(c.Request.Form["sn"], "")
	ts    := strings.Join(c.Request.Form["ts"], "")

	// 验证来源
	value, ok := config.ApiAuthConfig[ak]
	if ok {
		AppSecret = value["md5"]
	} else {
		return nil, errors.New("ak Error")
	}

	if debug == "1" {
		currentUnix := util.GetCurrentUnix()
		req.Set("ts", strconv.FormatInt(currentUnix, 10))
		res := map[string]string{
			"ts": strconv.FormatInt(currentUnix, 10),
			"sn": createSign(req),
		}
		return res, nil
	}

	// 验证过期时间
	timestamp := time.Now().Unix()
	exp, _    := strconv.ParseInt(config.AppSignExpiry, 10, 64)
	tsInt, _  := strconv.ParseInt(ts, 10, 64)
	if tsInt > timestamp || timestamp - tsInt >= exp {
		return nil, errors.New("ts Error")
	}

	// 验证签名
	if sn == "" || sn != createSign(req) {
		return nil, errors.New("sn Error")
	}

	return nil, nil
}

// 创建签名
func createSign(params url.Values) string {
	// 自定义 MD5 组合
	return util.MD5(AppSecret + createEncryptStr(params) + AppSecret)
}

func createEncryptStr(params url.Values) string {
	var key []string
	var str = ""
	for k := range params {
		if k != "sn" && k != "debug" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params.Get(key[i]))
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params.Get(key[i]))
		}
	}
	return str
}
```

## AES 对称加密

在使用前，咱们先了解下什么是对称加密？

对称加密就是使用同一个密钥即可以加密也可以解密，这种方法称为对称加密。

常用算法：DES、AES。

其中 AES 是 DES 的升级版，密钥长度更长，选择更多，也更灵活，安全性更高，速度更快，咱们直接上手 AES 加密。

**优点**

算法公开、计算量小、加密速度快、加密效率高。

**缺点**

发送方和接收方必须商定好密钥，然后使双方都能保存好密钥，密钥管理成为双方的负担。

**应用场景**

相对大一点的数据量或关键数据的加密。

#### 生成签名

首先，封装 Go 的 AesEncrypt 加密方法 和 AesDecrypt 解密方法。

```
// 加密 aes_128_cbc
func AesEncrypt (encryptStr string, key []byte, iv string) (string, error) {
	encryptBytes := []byte(encryptStr)
	block, err   := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes = pkcs5Padding(encryptBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// 解密
func AesDecrypt (decryptStr string, key []byte, iv string) (string, error) {
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted = pkcs5UnPadding(decrypted)
	return string(decrypted), nil
}

func pkcs5Padding (cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding (decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}
```

进行加密：

```
appKey     = "demo"
appSecret  = "xxx"
encryptStr = "param_1=xxx&param_2=xxx&ak="+appKey+"&ts=xxx"

sn = AesEncrypt(encryptStr, appSecret)
```

#### 验证签名

```
decryptStr = AesDecrypt(sn, app_secret)
```

将加密前的字符串与解密后的字符串做个对比。

相同，表示签名验证成功。

不同，表示签名验证失败。

#### 中间件 - 代码实现

```
var AppSecret string

// AES 对称加密
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		utilGin := util.Gin{Ctx: c}

		sign, err := verifySign(c)

		if sign != nil {
			utilGin.Response(-1, "Debug Sign", sign)
			c.Abort()
			return
		}

		if err != nil {
			utilGin.Response(-1, err.Error(), sign)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证签名
func verifySign(c *gin.Context) (map[string]string, error) {
	_ = c.Request.ParseForm()
	req   := c.Request.Form
	debug := strings.Join(c.Request.Form["debug"], "")
	ak    := strings.Join(c.Request.Form["ak"], "")
	sn    := strings.Join(c.Request.Form["sn"], "")
	ts    := strings.Join(c.Request.Form["ts"], "")

	// 验证来源
	value, ok := config.ApiAuthConfig[ak]
	if ok {
		AppSecret = value["aes"]
	} else {
		return nil, errors.New("ak Error")
	}

	if debug == "1" {
		currentUnix := util.GetCurrentUnix()
		req.Set("ts", strconv.FormatInt(currentUnix, 10))

		sn, err := createSign(req)
		if err != nil {
			return nil, errors.New("sn Exception")
		}

		res := map[string]string{
			"ts": strconv.FormatInt(currentUnix, 10),
			"sn": sn,
		}
		return res, nil
	}

	// 验证过期时间
	timestamp := time.Now().Unix()
	exp, _    := strconv.ParseInt(config.AppSignExpiry, 10, 64)
	tsInt, _  := strconv.ParseInt(ts, 10, 64)
	if tsInt > timestamp || timestamp - tsInt >= exp {
		return nil, errors.New("ts Error")
	}

	// 验证签名
	if sn == "" {
		return nil, errors.New("sn Error")
	}

	decryptStr, decryptErr := util.AesDecrypt(sn, []byte(AppSecret), AppSecret)
	if decryptErr != nil {
		return nil, errors.New(decryptErr.Error())
	}
	if decryptStr != createEncryptStr(req) {
		return nil, errors.New("sn Error")
	}
	return nil, nil
}

// 创建签名
func createSign(params url.Values) (string, error) {
	return util.AesEncrypt(createEncryptStr(params), []byte(AppSecret), AppSecret)
}

func createEncryptStr(params url.Values) string {
	var key []string
	var str = ""
	for k := range params {
		if k != "sn" && k != "debug" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params.Get(key[i]))
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params.Get(key[i]))
		}
	}
	return str
}
```

## RSA 非对称加密

和上面一样，在使用前，咱们先了解下什么是非对称加密？

非对称加密就是需要两个密钥来进行加密和解密，这两个秘钥分别是公钥（public key）和私钥（private key），这种方法称为非对称加密。

常用算法：RSA。

**优点**

与对称加密相比，安全性更好，加解密需要不同的密钥，公钥和私钥都可进行相互的加解密。

**缺点**

加密和解密花费时间长、速度慢，只适合对少量数据进行加密。

**应用场景**

适合于对安全性要求很高的场景，适合加密少量数据，比如支付数据、登录数据等。

#### 创建签名

首先，封装 Go 的 RsaPublicEncrypt 公钥加密方法 和 RsaPrivateDecrypt 解密方法。

```
// 公钥加密
func RsaPublicEncrypt(encryptStr string, path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件内容
	info, _ := file.Stat()
	buf := make([]byte,info.Size())
	file.Read(buf)

	// pem 解码
	block, _ := pem.Decode(buf)

	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//对明文进行加密
	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encryptStr))
	if err != nil {
		return "", err
	}

	//返回密文
	return base64.URLEncoding.EncodeToString(encryptedStr), nil
}

// 私钥解密
func RsaPrivateDecrypt(decryptStr string, path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 获取文件内容
	info, _ := file.Stat()
	buf := make([]byte,info.Size())
	file.Read(buf)

	// pem 解码
	block, _ := pem.Decode(buf)

	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)

	//对密文进行解密
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader,privateKey,decryptBytes)

	//返回明文
	return string(decrypted), nil
}

```

调用方 申请 公钥（public key），然后进行加密：

```
appKey     = "demo"
appSecret  = "公钥"
encryptStr = "param_1=xxx&param_2=xxx&ak="+appKey+"&ts=xxx"

sn = RsaPublicEncrypt(encryptStr, appSecret)
```

#### 验证签名

```
decryptStr = RsaPrivateDecrypt(sn, app_secret)
```

将加密前的字符串与解密后的字符串做个对比。

相同，表示签名验证成功。

不同，表示签名验证失败。

#### 中间件 - 代码实现

```
var AppSecret string

// RSA 非对称加密
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		utilGin := util.Gin{Ctx: c}

		sign, err := verifySign(c)

		if sign != nil {
			utilGin.Response(-1, "Debug Sign", sign)
			c.Abort()
			return
		}

		if err != nil {
			utilGin.Response(-1, err.Error(), sign)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证签名
func verifySign(c *gin.Context) (map[string]string, error) {
	_ = c.Request.ParseForm()
	req   := c.Request.Form
	debug := strings.Join(c.Request.Form["debug"], "")
	ak    := strings.Join(c.Request.Form["ak"], "")
	sn    := strings.Join(c.Request.Form["sn"], "")
	ts    := strings.Join(c.Request.Form["ts"], "")

	// 验证来源
	value, ok := config.ApiAuthConfig[ak]
	if ok {
		AppSecret = value["rsa"]
	} else {
		return nil, errors.New("ak Error")
	}

	if debug == "1" {
		currentUnix := util.GetCurrentUnix()
		req.Set("ts", strconv.FormatInt(currentUnix, 10))

		sn, err := createSign(req)
		if err != nil {
			return nil, errors.New("sn Exception")
		}

		res := map[string]string{
			"ts": strconv.FormatInt(currentUnix, 10),
			"sn": sn,
		}
		return res, nil
	}

	// 验证过期时间
	timestamp := time.Now().Unix()
	exp, _    := strconv.ParseInt(config.AppSignExpiry, 10, 64)
	tsInt, _  := strconv.ParseInt(ts, 10, 64)
	if tsInt > timestamp || timestamp - tsInt >= exp {
		return nil, errors.New("ts Error")
	}

	// 验证签名
	if sn == "" {
		return nil, errors.New("sn Error")
	}

	decryptStr, decryptErr := util.RsaPrivateDecrypt(sn, config.AppRsaPrivateFile)
	if decryptErr != nil {
		return nil, errors.New(decryptErr.Error())
	}
	if decryptStr != createEncryptStr(req) {
		return nil, errors.New("sn Error")
	}
	return nil, nil
}

// 创建签名
func createSign(params url.Values) (string, error) {
	return util.RsaPublicEncrypt(createEncryptStr(params), AppSecret)
}

func createEncryptStr(params url.Values) string {
	var key []string
	var str = ""
	for k := range params {
		if k != "sn" && k != "debug" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params.Get(key[i]))
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params.Get(key[i]))
		}
	}
	return str
}
```

## 如何调用？

与其他中间件调用方式一样，根据自己的需求自由选择。

比如，使用 MD5 组合：

```
.Use(sign_md5.SetUp())
```

使用 AES 对称加密：

```
.Use(sign_aes.SetUp())
```

使用 RSA 非对称加密：

```
.Use(sign_rsa.SetUp())
```

## 性能测试

既然 RSA 非对称加密，最安全，那么统一都使用它吧。

NO！NO！NO！绝对不行！

为什么我要激动，因为我以前遇到过这个坑呀，都是血泪的教训呀...

咱们挨个测试下性能：

#### MD5

```
func Md5Test(c *gin.Context) {
	startTime  := time.Now()
	appSecret  := "IgkibX71IEf382PT"
	encryptStr := "param_1=xxx&param_2=xxx&ak=xxx&ts=1111111111"
	count      := 1000000
	for i := 0; i < count; i++ {
		// 生成签名
		util.MD5(appSecret + encryptStr + appSecret)

		// 验证签名
		util.MD5(appSecret + encryptStr + appSecret)
	}
	utilGin := util.Gin{Ctx: c}
	utilGin.Response(1, fmt.Sprintf("%v次 - %v", count, time.Since(startTime)), nil)
}
```

模拟 一百万 次请求，大概执行时长在 1.1s ~ 1.2s 左右。

#### AES

```
func AesTest(c *gin.Context) {
	startTime  := time.Now()
	appSecret  := "IgkibX71IEf382PT"
	encryptStr := "param_1=xxx&param_2=xxx&ak=xxx&ts=1111111111"
	count      := 1000000
	for i := 0; i < count; i++ {
		// 生成签名
		sn, _ := util.AesEncrypt(encryptStr, []byte(appSecret), appSecret)

		// 验证签名
		util.AesDecrypt(sn, []byte(appSecret), appSecret)
	}
	utilGin := util.Gin{Ctx: c}
	utilGin.Response(1, fmt.Sprintf("%v次 - %v", count, time.Since(startTime)), nil)
}
```

模拟 一百万 次请求，大概执行时长在 1.8s ~ 1.9s 左右。

#### RSA

```
func RsaTest(c *gin.Context) {
	startTime  := time.Now()
	encryptStr := "param_1=xxx&param_2=xxx&ak=xxx&ts=1111111111"
	count      := 500
	for i := 0; i < count; i++ {
		// 生成签名
		sn, _ := util.RsaPublicEncrypt(encryptStr, "rsa/public.pem")

		// 验证签名
		util.RsaPrivateDecrypt(sn, "rsa/private.pem")
	}
	utilGin := util.Gin{Ctx: c}
	utilGin.Response(1, fmt.Sprintf("%v次 - %v", count, time.Since(startTime)), nil)
}
```

我不敢模拟 一百万 次请求，还不知道啥时候能搞定呢，咱们模拟 500 次试试。

模拟 500 次请求，大概执行时长在 1s 左右。

上面就是我本地的执行效果，大家可以质疑我的电脑性能差，封装的方法有问题...

你们也可以试试，看看性能差距是不是这么大。

## PHP 与 Go 加密方法如何互通？

我是写 PHP 的，生成签名的方法用 PHP 能实现吗？

肯定能呀！

我用 PHP 也实现了上面的 3 中方法，可能会有一些小调整，总体问题不大，相关 Demo 已上传到 github：

https://github.com/xinliangnote/Encrypt

好了，就到这了。

## 源码地址

https://github.com/xinliangnote/go-gin-api

## go-gin-api 系列文章

- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)
- [2. 规划项目目录和参数验证](https://mp.weixin.qq.com/s/11AuXptWGmL5QfiJArNLnA)
- [3. 路由中间件 - 日志记录](https://mp.weixin.qq.com/s/eTygPXnrYM2xfrRQyfn8Tg)
- [4. 路由中间件 - 捕获异常](https://mp.weixin.qq.com/s/SconDXB_x7Gan6T0Awdh9A)
- [5. 路由中间件 - Jaeger 链路追踪（理论篇）](https://mp.weixin.qq.com/s/28UBEsLOAHDv530ePilKQA)
- [6. 路由中间件 - Jaeger 链路追踪（实战篇）](https://mp.weixin.qq.com/s/Ea28475_UTNaM9RNfgPqJA)