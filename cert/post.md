# 如何创建和签名SSL / TLS证书？
我们将学习如何生成SSL / TLS证书签名请求（CSR）并由证书颁发机构（CA）对其进行签名。就本教程而言，我们不会将CSR提交给真实的CA。相反，我们将同时扮演两个角色：证书颁发机构和证书申请人。我们将使用openssl执行以下操作：

1. 生成证书颁发机构（CA）的私钥和自签名证书。
2. 为我们要使用TLS的Web服务器生成私钥及其配对的证书签名请求（CSR）。
3. 使用CA的私钥对我们的Web服务器的CSR进行签名，并取回已签名的证书。
4. 验证证书是否有效。

为了做所有这些事情，需要 `openssl` 或 `libressl` 已安装。

## 1. 生成CA的私钥和自签名证书
我们将使用的第一个命令是 [openssl req](https://man.openbsd.org/openssl.1#req)，它代表请求。此命令用于创建和处理证书签名请求。它也可以用于为CA创建自签名证书。
```
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem
```
该 `-x509` 选项用于告诉 openssl 输出自签名证书而不是证书请求。如果您不知道，[X509](https://en.wikipedia.org/wiki/X.509)只是公钥证书的标准格式。

该 `-newkey rsa:4096` 选项基本上告诉 openssl 同时创建新的RSA私钥（4096位）及其证书请求。当我们将其与 `-x509` 选项一起使用时，它将输出证书而不是证书请求。

下一个选项是 `-days 365`，它指定证书有效的天数。

然后我们使用 `-keyout` 选项告诉 openssl 将创建的私钥写入 `ca-key.pem` 文件。

最后 `-out` 选项告诉 openssl 将证书写入 `ca-cert.pem` 文件。

当我们运行此命令时，openssl将开始生成私钥。

![](./image/post-1.png)

执行命令后，将要求我们提供一个密码短语，它将用于在将私钥写入PEM文件之前对其进行加密。

为什么要加密？因为如果以某种方式黑客破坏了私钥文件，那么黑客将无法在不知道密码的情况下使用它执行任何操作解密。

接下来，openssl 将要求我们提供一些身份信息以生成证书：

- 国家区号，比如中国就是CN，美国就是US
- 州或省名称
- 城市名称
- 机构名称，随便写
- 单位名称，随便写
- 通用名称（或域名）
- 邮箱地址

![](./image/post-2.png)

这样，证书和私钥文件就成功生成了。

如果我们查看 `ca-key.pem` 私钥文件，我们可以看到它说是 `ENCRYPTED PRIVATE KEY`：
```
-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIJnDBOBgkqhkiG9w0BBQ0wQTApBgkqhkiG9w0BBQwwHAQI1hrj8iiNzhMCAggA
MAwGCCqGSIb3DQIJBQAwFAYIKoZIhvcNAwcECC2m96JG0fpdBIIJSE668ot/znRT
-----END ENCRYPTED PRIVATE KEY-----
```
对于 `ca-cert.pem`，该证书不是加密的，仅以base64编码，因为它只包含公钥，身份信息和签名，每个人都应该看到。

我们可以使用 [openssl x509](https://man.openbsd.org/openssl.1#x509) 命令显示此证书中编码的所有信息。此命令还可以用于签署证书请求，稍后我们会看到。

让我们执行：
```
openssl x509 -in ca-cert.pem -noout -text
```
在这里，我们使用 `-in` 选项来传递CA的证书文件。并且 `-noout` 选项告诉它不输出原始的base64编码值。相反，我们使用 `-text` 选项，因为我们希望以可读的文本格式显示它。

![](./image/post-3.png)

我要复制这些命令并将其保存到 `gen.sh` 脚本中。有了这个脚本，就可以自动化生成一组密钥和证书。

![](./image/post-4.png)

在进行第二步之前，我将向您展示另一种提供身份信息的方法，而无需像之前一样以交互方式输入身份信息。为此，我们必须将 `-subj`选项添加到 openssl req 命令：
```
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=CN/ST=GD/L=GZ/O=Ezzz/OU=Education/CN=*.ezzz.com/emailAddress=ezzz@gmail.com"
```
在此 `-subj` 字符串中：

- `/C=CN` 用于国家
- `/ST=GD` 适用于州或省
- `/L=GZ` 用于地区名称或城市
- `/O=Ezzz` 用于组织
- `/OU=Education` 用于组织单位
- `/CN=*.ezzz.com` 用于通用名称或域名
- `/emailAddress=ezzz@gmail.com` 用于电子邮件地址

现在，让我们在 `gen.sh` 脚本的顶部添加命令 `rm *.pem`（以删除所有pem文件），然后在终端中运行它。

![](./image/post-5.png)

仍然会提示输入密码，但是它不再要求提供身份信息，因为我们已经在 `-subj` 选项中提供了它们。

## 2. 生成Web服务器的私钥和CSR#
下一步是为我们的Web服务器生成私钥和CSR。

几乎与第一步中使用的命令相同。除此之外，这次我们不想对其自签名，因此我们应该删除 `-x509` 选项。 `-days` 选项也应删除，因为我们不是创建证书，而是创建CSR。
```
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=CN/ST=GD/L=GZ/O=Ezzz/OU=Education/CN=*.ezzz.com/emailAddress=ezzz@gmail.com"
```
输出密钥的名称为 `server-key.pem`，输出证书请求文件为 `server-req.pem`，并且 `-subj` 应包含我们的Web服务器信息。

这次，在 `server-req.pem` 文件中显示的是 CERTIFICATE REQUEST，而不是 `ca-cert.pem` 文件中的 CERTIFICATE。这是因为它不是以前的证书，而是证书签名请求。
```
-----BEGIN CERTIFICATE REQUEST-----
MIIE1jCCAr4CAQAwgZAxCzAJBgNVBAYTAkNOMRIwEAYDVQQIDAlHdWFuZ0Rvbmcx
EjAQBgNVBAcMCUd1YW5nWmhvdTEOMAwGA1UECgwFVmluY2kxEjAQBgNVBAsMCUVk
-----END CERTIFICATE REQUEST-----
```
现在，我们进入第3步并签署此请求。

## 3. 签署Web服务器的证书请求#
要对证书进行签名，我们将使用以前用来显示证书的相同的openssl x509命令。让我们打开终端并运行以下命令：
```
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
```
在此命令中，我们使用 `-req` 选项来告诉 openssl 我们将传递证书请求。我们使用 `-in` 选项，后跟请求文件的名称：`server-req.pem`。

接下来，我们使用 `-CA` 选项来传递CA的证书文件：ca-cert.pem。以及 `-CAkey` 选项，用于传入CA的私钥：ca-key.pem。

那么1个重要选项是 `-CAcreateserial`。基本上，CA必须确保它签署的每个证书都具有唯一的序列号。因此，使用此选项，将生成一个包含下一个序列号的文件（如果不存在）。

最后，我们使用 `-out` 选项指定输出证书写入的文件。

![](./image/post-6.png)

现在您可以在此处看到，因为CA的私钥已加密，所以openssl要求提供密码以对其进行解密，然后才能用于对证书进行签名。它是以防CA的私钥被黑的一种对策。

现在我们已经为我们的Web服务器获得了签名证书。让我们以纯文本格式将其打印出来。
```
openssl x509 -in server-cert.pem -noout -text
```
![](./image/post-7.png)

这是其唯一的序列号，我们还可以看到一个刚生成的 `ca-cert.srl` 文件，其中的`0D428A1B04465FD7EF6E671AAB0B16BEB9532BB0`序列号其实跟上图中的序列号是相同的，只是这里显示不一样。

默认情况下，证书有效期为30天。我们可以通过在签名命令中添加 `-days` 选项来更改它。
```
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
```
现在有效期已更改为60天。

证书可用于具有不同域名的多个网站。我们可以通过在签署证书申请时指定主题备用名称扩展名来实现。

`openssl x509` 命令的 `-extfile` 选项允许我们声明包含扩展名的文件。我们可以在此页面中中看到配置文件的格式。

我们可以使用一些东西作为替代名称，例如电子邮件，DNS或IP。我将创建一个新文件`server-ext.cnf` 写入以下内容：
```
subjectAltName=DNS:*.ezzz.com,DNS:*.ezzz.org,IP:0.0.0.0
```
在这里，我将DNS设置为多个域名：`*.ezzz.com`和`*.ezzz.org`。我还将IP设置为0.0.0.0，这将使用于我们在localhost上开发。

现在，在证书签名命令中，添加`-extfile`选项并传入扩展名配置文件的名称：
```
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
```
让我们为我们的Web服务器签名证书以纯文本格式将其打印出来。
```
openssl x509 -in server-cert.pem -noout -text
```
可以看到，证书文件具有了新的扩展名部分，其中包含我们选择的所有主题替代名称：

![](./image/post-8.png)

看起来我们的自动化脚本已经准备就绪，当然还是必须输入很多次密码来保护私钥。

如果我们只想使用它进行开发和测试，我们可以告诉 openssl 不要加密私钥，这样它就不会要求我们输入密码。

为此，我们将 `-nodes` 选项添加到` openssl req` 命令中，如下所示：

> gen.sh
```
rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=CN/ST=GD/L=DZ/O=Ezzz/OU=Education/CN=*.ezzz.com/emailAddress=ezzz@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=CN/ST=GD/L=DZ/O=Ezzz/OU=Education/CN=*.ezzz.com/emailAddress=ezzz@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text
```
现在，如果我们执行`gen.sh`，它将不再要求输入密码。如果我们看查看`ca-key.pem`私钥文件，它显示为私钥，而不是像以前那样显示为加密的私钥。
```
-----BEGIN PRIVATE KEY-----
MIIJQgIBADANBgkqhkiG9w0BAQEFAASCCSwwggkoAgEAAoICAQDs1R+kzFl3Te7H
0jsHVHtlSnx5bdcBdZNJKjEehjoy5ZU8b9yvAzTbbpu7Hwbd2N99P0ckK3azU2yO
-----END PRIVATE KEY-----
```
## 4. 验证证书是否有效#
最后，将向您展示如何验证证书是否有效。 我们可以使用openssl verify命令来做到这一点：
```
openssl verify -CAfile ca-cert.pem server-cert.pem
```
我们只要传递受信任的CA的证书和要验证的证书。如果返回OK，则证书有效。