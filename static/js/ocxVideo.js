/*
    ocx播放短视频、回放、直播
    DrawRect（)四边形
    DrawPolygon() 多边形
    StopDraw() 停止画框
    ClearShape() 清除画框
 */

function getOcxViewer(name) {
    if (window.document[name]) {
        return window.document[name];
    }
    if (navigator.appName.indexOf("Microsoft Internet") == -1) {
        if (document.embeds && document.embeds[name])
            return document.embeds[name];
    }
    else {
        return document.getElementById(name);
    }
}

/*
    直播
    name: 控件id
    devtype：设备类型 hik海康 dahu大华
    devaddress：设备ip
    devport：设备端口
    username：设备账号
    password：设备密码
    channel：通道号
 */
function ocxPlay(name, devtype, devaddress, devport, username, password, channel){
    if(devaddress==""){
        alert("ip地址为空");
    }else if(devport==""){
        alert("端口号为空");
    }else if(username==""){
        alert("账号为空");
    }else if(password==""){
        alert("密码为空");
    }else if(channel==""){
        alert("通道号为空");
    }
    getOcxViewer(name).PushParameter("devtype", devtype);
    getOcxViewer(name).PushParameter("devaddress", devaddress);
    getOcxViewer(name).PushParameter("devport", devport);
    getOcxViewer(name).PushParameter("username", username);
    getOcxViewer(name).PushParameter("password", password);
    getOcxViewer(name).PushParameter("channel", channel);
    getOcxViewer(name).RealPlay();
}

function ocxStop(name) {
    getOcxViewer(name).StopRealPlay();
}

/*
    回放
 */
function ocxPlayBack(name, devtype, devaddress, devport, username, password, channel, begintime, endtime){
    getOcxViewer(name).PushParameter("devtype", devtype);
    getOcxViewer(name).PushParameter("devaddress", devaddress);
    getOcxViewer(name).PushParameter("devport", devport);
    getOcxViewer(name).PushParameter("username", username);
    getOcxViewer(name).PushParameter("password", password);
    getOcxViewer(name).PushParameter("channel", channel);
    getOcxViewer(name).PushParameter("begintime", begintime);
    getOcxViewer(name).PushParameter("endtime", endtime);
    getOcxViewer(name).PlayBack();
}

function ocxPlayBackStop(name) {
    getOcxViewer(name).StopPlayBack();
}

/*
    短视频
 */

function ocxPlayVideo(name, url) {
    console.log(url);
    getOcxViewer(name).PushParameter("url",url);
    getOcxViewer(name).PlayVod();
}

function ocxPlayVideoStop(name) {
    getOcxViewer(name).StopPlayVod();
}
// 矩形
function ocxRectangle(name) {
    getOcxViewer(name).DrawRect1();
}
// 四边形
function ocxDrawRect(name) {
    getOcxViewer(name).DrawRect();
}
// 多边形
function ocxDrawPolygon(name) {
    getOcxViewer(name).DrawPolygon();
}
// 停止画框
function ocxStopDraw(name) {
    getOcxViewer(name).StopDraw();
}
// 清除信息
function ocxClearShape(name) {
    getOcxViewer(name).ClearShape();
}

// 画框
function ocxPushMultiShape(name,shapeString) {
    // var shapeString = '[[67,51],[68,204],[252,199],[275,64],[180,17],[67,51]];[[577,47],[422,172],[634,203],[836,152],[796,38],[577,47]];[[724,241],[527,317],[720,376],[922,325],[899,219],[724,241]];';
    getOcxViewer(name).PushMultiShape(shapeString);
}

// 获得视频宽度和高度
function ocxVideoSize(name) {
    var width = getOcxViewer(name).VideoWidth();
    var height = getOcxViewer(name).VideoHeight();
    return [width,height];
}