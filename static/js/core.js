$(function() {

	$.getJSON("/static/json/china.json").done(function(response) {

		jsMap.config("#map-01", response);

		jsMap.config("#map-02", response, {
			width: 900,
			height: 600,
			areaName: {
				show: true
			},
			clickCallback: function(id, name) {
				var adcode = "";
				switch(name) {
					case "北京":
						adcode = "110000";
						break;
					case "黑龙江":
						adcode = "230000";
						break;
					case "吉林":
						adcode = "220000";
						break;
					case "辽宁":
						adcode = "210000";
						break;
					case "河北":
						adcode = "130000";
						break;
					case "内蒙古":
						adcode = "150000";
						break;
					case "天津":
						adcode = "120000";
						break;
					case "河南":
						adcode = "410000";
						break;
					case "山东":
						adcode = "370000";
						break;
					case "山西":
						adcode = "140000";
						break;
					case "江苏":
						adcode = "320000";
						break;
					case "安徽":
						adcode = "340000";
						break;
					case "湖北":
						adcode = "420000";
						break;
					case "浙江":
						adcode = "330000";
						break;
					case "湖南":
						adcode = "430000";
						break;
					case "江西":
						adcode = "360000";
						break;
					case "福建":
						adcode = "350000";
						break;
					case "陕西":
						adcode = "610000";
						break;
					case "宁夏":
						adcode = "640000";
						break;
					case "甘肃":
						adcode = "620000";
						break;
					case "重庆":
						adcode = "500000";
						break;
					case "四川":
						adcode = "510000";
						break;
					case "贵州":
						adcode = "520000";
						break;
					case "云南":
						adcode = "530000";
						break;
					case "广西":
						adcode = "450000";
						break;
					case "广东":
						adcode = "440000";
						break;
					case "台湾":
						adcode = "710000";
						break;
					case "海南":
						adcode = "460000";
						break;
					case "青海":
						adcode = "630000";
						break;
					case "新疆":
						adcode = "650000";
						break;
					case "西藏":
						adcode = "540000";
						break;
					case "上海":
						adcode = "310000";
						break;
					case "香港":
						adcode = "810000";
						break;
					case "澳门":
						adcode = "820000";
						break;
					case "南海诸岛":
						adcode = "460000";
						break;

				}
				window.location="provincialMap.html?adcode="+adcode;
			}
		});

		jsMap.config("#map-03", response, {
			multiple: true
		});
		$("#get-multiple-1").on("click", function() {
			console.log(jsMap.multipleValue("#map-03"));
		})
		$("#get-multiple-2").on("click", function() {
			console.log(jsMap.multipleValue("#map-03", {
				type: "object"
			}));
		})

		jsMap.config("#map-04", response, {
			stroke: {
				width: 2,
				color: "#000"
			}
		});

		jsMap.config("#map-05", response, {
			fill: {
				basicColor: "#259200",
				hoverColor: "#57cb00",
				clickColor: "#2e6f18"
			}
		});

		jsMap.config("#map-06", response, {
			fill: {
				basicColor: {
					heilongjiang: "#ff5900",
					jilin: "#19bb00",
					liaoning: "#6800ff"
				},
				hoverColor: {
					heilongjiang: "#ff8c4e",
					jilin: "#1fe000",
					liaoning: "#954dff"
				},
				clickColor: {
					heilongjiang: "#c94600",
					jilin: "#159a00",
					liaoning: "#5200c9"
				}
			}
		});

		jsMap.config("#map-07", response, {
			disabled: {
				name: ["heilongjiang", "jilin", "liaoning"]
			}
		});

		jsMap.config("#map-08", response, {
			disabled: {
				name: ["heilongjiang", "jilin", "liaoning"],
				except: true
			}
		});

		jsMap.config("#map-09", response, {
			tip: function(id, name) {
				return '<div style="background:#eee;padding:15px;"><p>id: ' + id + '</p><p>name: ' + name + '</p></div>';
			}
		});

		var $hoverCallback = $("#hover-callback");
		jsMap.config("#map-10", response, {
			hoverCallback: function(id, name) {
				$hoverCallback.text(id + " --- " + name);
			}
		});

		var $clickCallback = $("#click-callback");

	})

})