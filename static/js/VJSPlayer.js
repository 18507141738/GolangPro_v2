

function VJSPlayer(parOptions){
	parOptions = parOptions || {};
	var pthis = this;
	this.width = parOptions.width || 704;
	this.height = parOptions.height || 576;
	var yuvlen = this.width * this.height*3/2;
	this.serverHost = "localhost:3080"; 
	this.netConfig =  {};			
	this.realSessid = "";
	this.realFramesUrl = "";
	this.recordSessid = "";
	this.recordFramesUrl = "";
	this.downloadSessid = "";
	this.canvas = parOptions.canvas;
	this.glcanvas =  new YUVCanvas({     
							webgl : true,
							type: 'yuv420',
							canvas: parOptions.canvas,
							width: this.width,
							height: this.height});
	
	this.resq = {};
	this.seqCtx = { 'realSeq':0,'recordSeq':0,'downSeq':0};
	this.seqNum = 0;
	this.reconnTimes = 0;
	this.duration = 0;
	this.pos = 0;
	this.playbackStatus = null;
	this.timer = setInterval(_timeout(pthis),1000);
	this.lastTime = new Date().getTime();
}

VJSPlayer.prototype.fullscreen = function(width,height){
	if(this.canvas.webkitRequestFullscreen)
		this.canvas.webkitRequestFullscreen();
	else if(this.canvas.exitFullscreen)
		this.canvas.exitFullscreen();
	else if(this.canvas.msRequestFullscreen)
		this.canvas.msRequestFullscreen();
	else
		this.canvas.requestFullscreen();
	return ;
}


VJSPlayer.prototype.width = function(){
	return this.width;
}

VJSPlayer.prototype.height = function(){
	return this.height;
}

//改变播放窗口大小
VJSPlayer.prototype.resize = function(width,height){
	var url = "http://"+this.serverHost+"/webcontrol/real_resize";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.width = width;
	resq.height = height;
	resq.sessid = this.realSessid;
	xhr.timeout = 5000;
	xhr.responseType = "json";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				console.log("resize (" + width + "," + height);
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
}

//VJSPlayer.prototype.timeout = function (p){
	
function timeout(pthis){
	//console.log('timeout'+ pthis.width + pthis.serverHost);
	var curTime = new Date().getTime();
	let diffTime = 5000 + (pthis.reconnTimes*pthis.reconnTimes)*1000;
	if(diffTime > 90000)
		diffTime = 90000;
	if(curTime - pthis.lastTime  > diffTime && pthis.realSessid != ""){
		pthis.realPlay(pthis.resq);
		pthis.reconnTimes++;
	}
	
	if(pthis.playbackStatus && pthis.recordSessid != ""){
		pthis.playbackStatus(pthis.recordSessid,pthis.duration,pthis.pos,"",0);
	}
}

function _timeout(param){
	return function(){
		timeout(param);
	}
}

//实时流播放
VJSPlayer.prototype.realPlay = function(resq){
	
	if(this.realSessid != "")
	{
		this.realStop(this.realSessid);
		this.realSessid = "";
	}
	if(this.recordSessid != "")
	{
		this.recordStop(this.recordSessid);
		this.recordSessid= "";
	}
	if(this.downloadSessid != "")
	{
		this.downloadStop(this.downloadSessid);
		this.downloadSessid = "";
	}
	
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/real_play";
	var xhr = new XMLHttpRequest();
	this.seqNum++;
	this.seqCtx.realSeq = this.seqNum;
	this.lastTime = new Date().getTime();
	this.resq = resq;
	resq.seq = this.seqNum;
	xhr.timeout = 5000;
	xhr.responseType = "json";//"blob","arraybuffer";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				let seq = this.response.data.seq;
				let realFramesUrl = this.response.data.url;
				let realSessid = this.response.data.sessid;
				
				if(pthis.seqNum != seq)
					pthis.realStop(realSessid);
				else{
					pthis.realFramesUrl = this.response.data.url;
					pthis.realSessid = this.response.data.sessid;
					pthis.realFrames(pthis.realFramesUrl);
				}
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
};

VJSPlayer.prototype.realFrames = function(url){
	var pthis = this;
	var frame_type = pthis.glcanvas.isWebGL() ? 0:1; //0 YUV 1 RGB
	var maxFrameLen = this.width* this.height*3 + 1024;
	var frameBuffer =  new Uint8Array(maxFrameLen);
    var buffLength = 0;
	url += "&width=" + this.width + "&height="+this.height + "&frame_type="+frame_type;
	fetch(url).then((response) => {
			  const reader = response.body.getReader();
			  const stream = response.body;
				  //下面的函数处理每个数据块
				  function push() {
					//"done"是一个布尔型，"value"是一个Unit8Array
					reader.read().then(({ done, value }) => {
							
					
					    //console.log(++cc+': ' +    value.length  );
						
						frameBuffer.set(value,buffLength);
						buffLength += value.length;
						
						var pkglen =  frameBuffer[3] << 24 
									| frameBuffer[2] << 16 
									| frameBuffer[1] << 8 
									| frameBuffer[0] << 0;
									
						if(pkglen > maxFrameLen){
							let tmp = frameBuffer.subarray(0,buffLength);
							maxFrameLen =  pkglen*2;
							frameBuffer = new Uint8Array(maxFrameLen);
							frameBuffer.set(tmp,0);
						}
							
										
						if(buffLength >= pkglen){
							
							var width  =  frameBuffer[187] << 24 
										| frameBuffer[186] << 16 
										| frameBuffer[185] << 8 
										| frameBuffer[184] << 0;
							var height  = frameBuffer[191] << 24 
										| frameBuffer[190] << 16 
										| frameBuffer[189] << 8 
										| frameBuffer[188] << 0;
							pthis.duration =  frameBuffer[271] << 24 
										| frameBuffer[270] << 16 
										| frameBuffer[269] << 8 
										| frameBuffer[268] << 0;
							pthis.pos =  frameBuffer[275] << 24 
										| frameBuffer[274] << 16 
										| frameBuffer[273] << 8 
										| frameBuffer[272] << 0;
										
							if(width <= 0 || width > 5000 || height <= 0 || height > 5000){
								controller.close();
								return;
							}
							
							if(width != pthis.width || height != pthis.height){
								 pthis.glcanvas =  new YUVCanvas({     
														webgl : true,
														type: 'yuv420',
														canvas: pthis.canvas,
														width: width,
														height: height});
								pthis.width = width; 
								pthis.height = height;
							}
							
							var hdrlen = 1024;
							var yuvlen = width * height*3/2;
							var ylen = width * height;
							var uvlen = (width / 2) * (height / 2);
							var framelen = pkglen;// yuvlen + hdrlen;
							pthis.lastTime = new Date().getTime();
							pthis.reconnTimes = 0;
						
							if(frame_type)
								pthis.glcanvas.drawNextOutputPicture(width,height,null, frameBuffer.subarray(hdrlen, pkglen) )
							else
								pthis.glcanvas.drawNextOutputPicture(width,height,null,{
										yData: frameBuffer.subarray(hdrlen + 0, hdrlen + ylen),
										uData: frameBuffer.subarray(hdrlen+ ylen, hdrlen + ylen + uvlen),
										vData: frameBuffer.subarray(hdrlen+ ylen + uvlen, hdrlen + ylen + uvlen + uvlen)});
							
							let remain = frameBuffer.subarray(framelen,buffLength);
							frameBuffer.set(remain,buffLength-framelen);//frameBuffer.copyWithIn(0,yuvlen,bufflen);
							buffLength -= framelen;
						}
						
					  //判断是否还有可读的数据？
					  if (done) {
						//告诉浏览器已经结束数据发送。
						controller.close();
						return;
					  }
					  //取得数据并将它通过controller发送给浏览器。
					  //controller.enqueue(value);
					  push();
					}).catch(function(error) {
							console.log('reader.read() failed', error)
						  });
				  
				  }
				  push();
			  return new Response(stream, { headers: { "Content-Type": "application/octet-stream" } });
			}).catch(function(error) {
					console.log('request failed', error)
				  });
}

//播放时截图
VJSPlayer.prototype.realSnap = function(picpath){
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/real_snap";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.sessid = this.realSessid;
	resq.file_path = picpath;
	xhr.timeout = 5000;
	xhr.responseType = "json";//"blob","arraybuffer";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				console.log('real_snap successful');
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	
	return ;
};

//停止实时流播放
VJSPlayer.prototype.realStop = function(){
	if(this.realSessid == "")
		return ;
	
	var url = "http://"+ this.serverHost+"/webcontrol/real_stop";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.sessid = this.realSessid;
	xhr.timeout = 5000;
	xhr.responseType = "json";
	xhr.open('POST', url, true);
	xhr.onload = function() {
		if (this.status == 200) {
			console.log("real stop finished !");
		 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	this.realSessid = "";
	return ;
}

//录像流播放
VJSPlayer.prototype.recordPlay = function(resq){
	
	if(this.realSessid != "")
	{
		this.realStop(this.realSessid);
		this.realSessid = "";
	}
	if(this.recordSessid != "")
	{
		this.recordStop(this.recordSessid);
		this.recordSessid= "";
	}
	if(this.downloadSessid != "")
	{
		this.downloadStop(this.downloadSessid);
		this.downloadSessid = "";
	}
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/record_play";
	var xhr = new XMLHttpRequest();
	this.seqNum++;
	this.seqCtx.recordSeq = this.seqNum;
	resq.seq = this.seqNum;
	xhr.timeout = 5000;
	xhr.responseType = "json";//"blob","arraybuffer";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				let seq = this.response.data.seq;
				let recordFramesUrl = this.response.data.url;
				let recordSessid = this.response.data.sessid;
				
				if(pthis.seqNum != seq){
					recordStop(recordSessid);
				}
				else{
					pthis.recordFramesUrl = this.response.data.url;
					pthis.recordSessid = this.response.data.sessid;
					pthis.realFrames(pthis.recordFramesUrl);
				}
					
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
};

//截图
VJSPlayer.prototype.recordSnap = function(picpath){
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/record_snap";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.sessid = this.recordSessid;
	resq.file_path = picpath;
	xhr.timeout = 5000;
	xhr.responseType = "json";//"blob","arraybuffer";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				console.log('record_snap successful');
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
};

//跳到指定时间播放
VJSPlayer.prototype.recordSetpos= function(pos){
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/record_set_pos";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.sessid = this.recordSessid;
	resq.pos = pos;
	xhr.timeout = 5000;
	xhr.responseType = "json";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				console.log('record_set_pos successful');
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
};

//停止播放
VJSPlayer.prototype.recordStop = function(){
	if(this.recordSessid == "")
		return ;
	
	var url = "http://"+ this.serverHost+"/webcontrol/record_stop";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.sessid = this.recordSessid;
	xhr.timeout = 5000;
	xhr.responseType = "json";
	xhr.open('POST', url, true);
	xhr.onload = function() {
		if (this.status == 200) {
			console.log("real stop finished !");
		 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	this.recordSessid  = "";
	return ;
}

//文件下载
VJSPlayer.prototype.downloadStart = function(resq){
	
	if(this.realSessid != "")
	{
		this.realStop(this.realSessid);
		this.realSessid = "";
	}
	if(this.recordSessid != "")
	{
		this.recordStop(this.recordSessid);
		this.recordSessid= "";
	}
	if(this.downloadSessid != "")
	{
		this.downloadStop(this.downloadSessid);
		this.downloadSessid = "";
	}
	
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/download_start";
	var xhr = new XMLHttpRequest();
	this.seqNum++;
	this.seqCtx.downSeq = this.seqNum;
	resq.seq = this.seqNum;
	xhr.timeout = 5000;
	xhr.responseType = "json";//"blob","arraybuffer";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				var seq = this.response.data.seq;
				pthis.downloadSessid = this.response.data.sessid;
				if(pthis.seqNum != seq)
					downloadStop(pthis.downloadSessid);
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
};

//下载状态
VJSPlayer.prototype.downloadStatus = function(){
	var pthis = this;
	var url = "http://"+this.serverHost+"/webcontrol/download_status";
	var xhr = new XMLHttpRequest();
	xhr.timeout = 5000;
	xhr.responseType = "json";//"blob","arraybuffer";
	xhr.open('POST', url, true);
	xhr.onload = function() {
			if (this.status == 200) {
				var status = this.response.data.status;
				var progress = this.response.data.progress;
				console.log("download status: "+status + " progress: " + progress);
			 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	return ;
};

//停止下载
VJSPlayer.prototype.downloadStop = function(){
	if( this.downloadSessid == "")
		return ;
	
	var url = "http://"+ this.serverHost+"/webcontrol/download_stop";
	var xhr = new XMLHttpRequest();
	var resq = {};
	resq.sessid = this.downloadSessid;
	xhr.timeout = 5000;
	xhr.responseType = "json";
	xhr.open('POST', url, true);
	xhr.onload = function() {
		if (this.status == 200) {
			console.log("download stop finished !");
		 }
	 }
	xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(resq));	
	this.downloadSessid = "";
	return ;
};
