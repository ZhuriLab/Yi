package runner

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

var ConfigFileName = "config.yaml"

// 默认配置文件,  todo 注: Codeql 不支持指定文件夹来运行规则
var defaultYamlByte = []byte(`
go_ql:
  # sql 注入
  - go/ql/src/Security/CWE-089/SqlInjection.ql
  # LDAP 注入
  - go/ql/src/experimental/CWE-090/LDAPInjection.ql
  # 登录认证绕过
  - go/ql/src/experimental/CWE-285/PamAuthBypass.ql
  # 用户控制绕过
  - go/ql/src/experimental/CWE-807/SensitiveConditionBypass.ql
  # JWT 硬编码
  - go/ql/src/experimental/CWE-321/HardcodedKeys.ql
  # SSRF
  - go/ql/src/experimental/CWE-918/SSRF.ql
  # 路径遍历
  - go/ql/src/Security/CWE-022/TaintedPath.ql
  # 解压造成的路径穿越
  - go/ql/src/Security/CWE-022/UnsafeUnzipSymlink.ql
  - go/ql/src/Security/CWE-022/ZipSlip.ql
  # 命令执行
  - go/ql/src/Security/CWE-078/CommandInjection.ql
  - go/ql/src/Security/CWE-078/StoredCommand.ql
  # xss
  - go/ql/src/Security/CWE-079/ReflectedXss.ql
  - go/ql/src/Security/CWE-079/StoredXss.ql
  # debug 、堆栈信息泄露
  - go/ql/src/Security/CWE-209/StackTraceExposure.ql
  # 重定向
  - go/ql/src/Security/CWE-601/BadRedirectCheck.ql
  - go/ql/src/Security/CWE-601/OpenUrlRedirect.ql
  # Xpath 注入
  - go/ql/src/Security/CWE-643/XPathInjection.ql
  # 硬编码认证信息
  - go/ql/src/Security/CWE-798/HardcodedCredentials.ql	
  # 重定向
  - go/ql/src/myRules/UrlRedirect.ql
  # 任意文件读取
  - go/ql/src/myRules/ReadFile.ql
java_ql:
  # 路径问题
  - java/ql/src/Security/CWE/CWE-022/TaintedPath.ql
  - java/ql/src/Security/CWE/CWE-022/TaintedPathLocal.ql
  - java/ql/src/Security/CWE/CWE-022/ZipSlip.ql
  # JNDI 
  - java/ql/src/Security/CWE/CWE-074/JndiInjection.ql
  - java/ql/src/Security/CWE/CWE-074/XsltInjection.ql
  # 命令执行
  - java/ql/src/Security/CWE/CWE-078/ExecRelative.ql
  - java/ql/src/Security/CWE/CWE-078/ExecTainted.ql
  - java/ql/src/Security/CWE/CWE-078/ExecTaintedLocal.ql
  - java/ql/src/Security/CWE/CWE-078/ExecUnescaped.ql
  # xss
  - java/ql/src/Security/CWE/CWE-079/XSS.ql
  - java/ql/src/Security/CWE/CWE-079/XSSLocal.ql
  # sql 注入
  - java/ql/src/Security/CWE/CWE-089/SqlTainted.ql
  - java/ql/src/Security/CWE/CWE-089/SqlTaintedLocal.ql
  - java/ql/src/Security/CWE/CWE-089/SqlUnescaped.ql
  # LDAP
  - java/ql/src/Security/CWE/CWE-090/LdapInjection.ql
  # injection
  - java/ql/src/Security/CWE/CWE-094/GroovyInjection.ql
  - java/ql/src/Security/CWE/CWE-094/InsecureBeanValidation.ql
  - java/ql/src/Security/CWE/CWE-094/JexlInjection.ql
  - java/ql/src/Security/CWE/CWE-094/MvelInjection.ql
  - java/ql/src/Security/CWE/CWE-094/SpelInjection.ql
  - java/ql/src/Security/CWE/CWE-094/TemplateInjection.ql
  # debug 、堆栈信息泄露
  - java/ql/src/Security/CWE/CWE-209/StackTraceExposure.ql
  # CVE-2019-16303
  - java/ql/src/Security/CWE/CWE-338/JHipsterGeneratedPRNG.ql
  # JWT 
  - java/ql/src/Security/CWE/CWE-347/MissingJWTSignatureCheck.ql
  # 反序列化
  - java/ql/src/Security/CWE/CWE-502/UnsafeDeserialization.ql
  # 敏感信息写入日志
  - java/ql/src/Security/CWE/CWE-532/SensitiveInfoLog.ql
  # XXE
  - java/ql/src/Security/CWE/CWE-611/XXE.ql
  # XPATH
  - java/ql/src/Security/CWE/CWE-643/XPathInjection.ql
  # Dos
  - java/ql/src/Security/CWE/CWE-730/ReDoS.ql
  # 硬编码
  - java/ql/src/Security/CWE/CWE-798/HardcodedCredentialsApiCall.ql
  - java/ql/src/Security/CWE/CWE-798/HardcodedCredentialsComparison.ql
  - java/ql/src/Security/CWE/CWE-798/HardcodedCredentialsSourceCall.ql
  - java/ql/src/Security/CWE/CWE-798/HardcodedPasswordField.ql
  # 用户权限绕过
  - java/ql/src/Security/CWE/CWE-807/ConditionalBypass.ql
  - java/ql/src/Security/CWE/CWE-807/TaintedPermissionsCheck.ql
  # OGNL
  - java/ql/src/Security/CWE/CWE-917/OgnlInjection.ql
  # spring boot
  - java/ql/src/experimental/Security/CWE/CWE-016/InsecureSpringActuatorConfig.ql
  - java/ql/src/experimental/Security/CWE/CWE-016/SpringBootActuators.ql
  # log4j
  - java/ql/src/experimental/Security/CWE/CWE-020/Log4jJndiInjection.ql
  # file
  - java/ql/src/experimental/Security/CWE/CWE-073/FilePathInjection.ql
  # 命令执行
  - java/ql/src/experimental/Security/CWE/CWE-078/ExecTainted.ql
  # MyBatis sql
  - java/ql/src/experimental/Security/CWE/CWE-089/MyBatisAnnotationSqlInjection.ql
  - java/ql/src/experimental/Security/CWE/CWE-089/MyBatisMapperXmlSqlInjection.ql
  # injection
  - java/ql/src/experimental/Security/CWE/CWE-094/BeanShellInjection.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/InsecureDexLoading.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/JakartaExpressionInjection.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/JShellInjection.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/JythonInjection.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/ScriptInjection.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/SpringImplicitViewManipulation.ql
  - java/ql/src/experimental/Security/CWE/CWE-094/SpringViewManipulation.ql
  # JWT 
  - java/ql/src/experimental/Security/CWE/CWE-321/HardcodedJwtKey.ql
  # Jsonp
  - java/ql/src/experimental/Security/CWE/CWE-352/JsonpInjection.ql
  # 反射
  - java/ql/src/experimental/Security/CWE/CWE-470/UnsafeReflection.ql
  # 反序列化
  - java/ql/src/experimental/Security/CWE/CWE-502/UnsafeDeserializationRmi.ql
  - java/ql/src/experimental/Security/CWE/CWE-502/UnsafeSpringExporterInConfigurationClass.ql
  - java/ql/src/experimental/Security/CWE/CWE-502/UnsafeSpringExporterInXMLConfiguration.ql
  # 目录
  - java/ql/src/experimental/Security/CWE/CWE-548/InsecureDirectoryConfig.ql
  - java/ql/src/experimental/Security/CWE/CWE-552/UnsafeUrlForward.ql
  # debug 
  - java/ql/src/experimental/Security/CWE/CWE-600/UncaughtServletException.ql
  # 重定向
  - java/ql/src/experimental/Security/CWE/CWE-601/SpringUrlRedirect.ql
  # XXE
  - java/ql/src/experimental/Security/CWE/CWE-611/XXE.ql
  - java/ql/src/experimental/Security/CWE/CWE-611/XXELocal.ql
  # 正则导致的权限绕过问题
  - java/ql/src/experimental/Security/CWE/CWE-625/PermissiveDotRegex.ql
  # 注入
  - java/ql/src/experimental/Security/CWE/CWE-652/XQueryInjection.ql
  # RMI
  - java/ql/src/experimental/Security/CWE/CWE-665/InsecureRmiJmxEnvironmentConfiguration.ql
  # Dos
  - java/ql/src/experimental/Security/CWE/CWE-755/NFEAndroidDoS.ql
`)
