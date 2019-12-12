module.exports = {
    runtimeCompiler: true,
    devServer: {
        proxy: {
            // '/apis': 'https://127.0.0.1:8081',
            '/apis': {    //将www.exaple.com印射为/apis
                target: 'http://127.0.0.1:8081',  // 接口域名
                secure: false,  // 如果是https接口，需要配置这个参数
                changeOrigin: true,  //是否跨域
                // pathRewrite: {
                //     '^/apis': ''   //需要rewrite的,
                // }
            },
            '/img': {    //将www.exaple.com印射为/apis
                target: 'http://127.0.0.1:8081',  // 接口域名
                secure: false,  // 如果是https接口，需要配置这个参数
                changeOrigin: true,  //是否跨域
                // pathRewrite: {
                //     '^/apis': ''   //需要rewrite的,
                // }
            }
        }
    }
}