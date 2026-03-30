setTimeout(function(){
    var userInfo = localStorage.getItem('userInfo');
    if(userInfo){
        var user = JSON.parse(userInfo);
        var email = user.Email || user.email;
        if(email){
            var url = '/api/esg/list?userEmail=' + encodeURIComponent(email);
            console.log('[Dashboard] 开始加载文件列表...');
            fetch(url).then(function(r){return r.json()}).then(function(r){
                if(r.code === 200){
                    console.log('[Dashboard] 成功加载 ' + r.data.length + ' 个文件');
                    if(typeof displayFileList === 'function'){
                        displayFileList(r.data);
                    }
                }else{
                    console.error('[Dashboard] 加载失败:', r.message);
                }
            }).catch(function(e){
                console.error('[Dashboard] 异常:', e);
            });
        }else{
            console.log('[Dashboard] 用户信息中没有邮箱');
        }
    }else{
        console.log('[Dashboard] 未找到用户信息，可能未登录');
    }
}, 1000);
