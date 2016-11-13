package html

const Template = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Degrees Of Separation</title>
		<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" crossorigin="anonymous">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" crossorigin="anonymous"></script>
		<script type="text/javascript">
		  	function checkDoS(){
			    //alert("Degrees Of Separation");
			    $("#result-div").show();
			    $("#loader").show();
			    $.post('checkDoS',{"actor1":$("#actor1").val(),"actor2":$("#actor2").val()},function(data){
			    	//alert(data);			    	
			    	$("#loader").hide();
			    	$("#result-data").html("<h2>Degrees Of Separation</h2>"+data);
			    });
			}
		</script>
		<style type="text/css">
			.grey-background {
			    background-color: #2f2f2f;
			    color: #f2f0e1;
			    overflow: hidden;
			}
			h1{
				text-align: center;
			}
			.well{
				background-color: #ebebe8;
    			color: #000;
    			width: 455px;
    			height: 290px;
			}
			.result {
			    background-color: #ebebe8;
			    color: #000;
			    text-align: left !important;
			    width: 65%;
			    min-height: 150px;
			    border-radius: 5px;
			    padding: 5px 15px;
			    overflow: auto;
			}
			.result-data {
				font-size: 18px;
			}
			form{
				text-align: left !important;
			}
			label {
			    font-size: 18px;
			}
			.txt_box {
			    border-style: solid;
			    border-width: 1px;
			    border-color: rgb(204, 204, 204);
			    border-radius: 5px;
			    background-color: rgb(255, 255, 255);
			    box-shadow: 0.5px 0.866px 2px 0px rgba(0, 0, 0, 0.063);
			    width: 350px;
			    height: 50px;
			    padding-left: 15px;
			    margin-bottom: 15px;
			}
			.bos_btn:hover, .bos_btn:active {
			    background-color: rgba(100, 52, 149,0.802);
			    border-color: rgb(100, 52, 149);
			    color: #FFF;
			    text-decoration: none;
			    cursor: pointer;
			}
			.bos_btn {
			    font-size: 14px;
			    color: rgb(255, 255, 255);
			    font-weight: bold;
			    text-align: center;
			    border-radius: 5px;
			    background-color: #673399;
			    box-shadow: 0.5px 0.866px 2px 0px rgba(0, 0, 0, 0.063);
			    width: 175px;
			    height: 50px;
			    padding: 14px;
			    margin-right: 35px;
			}
		</style>
	</head>
	<body class="grey-background">
		<h1>Degrees Of Separation</h1>
		<div align="center">
			<div class="well">
				<form>
				  <label>First name:</label><br>
				  <input type="text" id="actor1" name="actor1" value="vijay" class="txt_box">
				  <br>
				  <label>Last name:</label><br>
				  <input type="text" id="actor2" name="actor2" value="ajith-kumar" class="txt_box">
				  <br><br>
				  <a onclick="checkDoS()" class="bos_btn pull-right">Check DoS!</a>
				</form>	
			</div>	
			<div class="result" id="result-div" style="display: none;">
				<div id="loader" align="center" style="padding: 8% 0px;display: none;color: #2f2f2f;">
					<i class="fa fa-refresh fa-spin fa-5x fa-fw"></i>				
				</div>				
				<div class="result-data" id="result-data">
				
				</div>
			</div>
		</div>		
	</body>
</html>`