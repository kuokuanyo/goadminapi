package refund

// 登入介面HTML語法
const refundTmpl = `{{define "refund_theme1"}}
    <!DOCTYPE html>
    <!--[if lt IE 7]>
    <html class="no-js lt-ie9 lt-ie8 lt-ie7">
    <![endif]-->
    <!--[if IE 7]>
    <html class="no-js lt-ie9 lt-ie8">
    <![endif]-->
    <!--[if IE 8]>
    <html class="no-js lt-ie9">
    <![endif]-->
    <!--[if gt IE 8]><!-->
    <html class="no-js">
    <!--<![endif]-->
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
		<title>退費</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link rel="stylesheet" href="{{link .CdnUrl .UrlPrefix "/assets/refund/dist/all.min.css"}}">

        <!--[if lt IE 9]>
        <script src="{{link .CdnUrl .UrlPrefix "/assets/refund/dist/respond.min.js"}}"></script>
        <![endif]-->

    </head>
    <body>

    <div class="container">
        <div class="row" style="margin-top: 80px;">
            <div class="col-md-4 col-md-offset-4">
                <form action="##" onsubmit="return false" method="post" id="refund-form" class="fh5co-form animate-box"
                      data-animate-effect="fadeIn">
					<h2>申請退費</h2>
                    <div class="form-group">
						<tr>
							<label for="reason" class="sr-only">退費原因</label>
							<td><input type="text" class="form-control" id="reason" autocomplete="off" required="required" placeholder="輸入退費原因" ></td>
						<tr/>
                    </div>
                    <div class="form-group">
                        <button class="btn btn-primary" onclick="submitData()">提交</button>
                    </div>
                </form>
            </div>
        </div>
        <div class="row" style="padding-top: 60px; clear: both;">
            <div class="col-md-12 text-center"></div>
        </div>
    </div>

    <div id="particles-js">
        <canvas class="particles-js-canvas-el" width="1606" height="1862" style="width: 100%; height: 100%;"></canvas>
    </div>

    <script src="{{link .CdnUrl .UrlPrefix "/assets/refund/dist/all.min.js"}}"></script>

    <script>
        function submitData() {
            $.ajax({
                dataType: 'json',
                type: 'POST',
                url: '{{.UrlPrefix}}/refund',
                async: 'true',
                data: {
                    'reason': $("#reason").val(),
                    'orderID': getQueryVariable("orderID")
                },
                success: function (data) {
					location.href = "http://www.cco.com.tw"
                    alert("成功提交退費申請並可以關閉此頁面，完成退費後會傳送訊息給您!\n");
                },
                error: function (data) {
					alert(data.responseJSON.msg);
                }
            });
		}
		
	    function getQueryVariable(variable){
			var query = window.location.search.substring(1);
			var vars = query.split("&");
			for (var i=0;i<vars.length;i++) {
					var pair = vars[i].split("=");
					if(pair[0] == variable){return pair[1];}
			}
       		return(false);
		}

		function getCharFromUtf8(str) {  
			var cstr = "";  
			var nOffset = 0;  
			if (str == "")  
			return "";  
				str = str.toLowerCase();  
				nOffset = str.indexOf("%e");  
			if (nOffset == -1)  
			return str;  
			while (nOffset != -1) {  
					cstr += str.substr(0, nOffset);  
					str = str.substr(nOffset, str.length - nOffset);  
			if (str == "" || str.length < 9)  
			return cstr;  
					cstr += utf8ToChar(str.substr(0, 9));  
					str = str.substr(9, str.length - 9);  
					nOffset = str.indexOf("%e");  
				}  
			return cstr + str;  
		} 

		function utf8ToChar(str) {  
			var iCode, iCode1, iCode2;  
				iCode = parseInt("0x" + str.substr(1, 2));  
				iCode1 = parseInt("0x" + str.substr(4, 2));  
				iCode2 = parseInt("0x" + str.substr(7, 2));  
			return String.fromCharCode(((iCode & 0x0F) << 12) | ((iCode1 & 0x3F) << 6) | (iCode2 & 0x3F));  
			} 

    </script>

    </body>
    </html>
{{end}}`
