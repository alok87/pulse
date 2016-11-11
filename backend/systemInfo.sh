#!/bin/bash 

export TZ=Asia/Calcutta

# Variables and declarations----------------------------------
FSEP="#"
f_header="./templates/_header.html"
report="./results/report.html"
f_config="./config"
dirNexus="http://nexus/service/local/repositories/15_OFS/content/pulse/static"
f_css_1="http://nexus/service/local/repositories/15_OFS/content/pulse/static/bootstrap.css"
f_css_2="http://nexus/service/local/repositories/15_OFS/content/pulse/static/bootswatch.min.css"
const_width=5
f_status="${HOME}/pulse/data/platform/dev/platform_oms/platform_oms.dat"
f_status_htm="${HOME}/pulse/data/platform/dev/platform_oms/platform_oms.html" 
NEXUS_USER=""
NEXUS_PASS=""
NEXUS_UPLOAD_URL="http://nexus/content/repositories/"


# Fucntions -------------------------------------------------------

function uploadNexus {
     curl -v -u ${NEXUS_USER}:${NEXUS_PASS} --upload-file $1 $2 > /dev/null 2>&1
     if [ $? -ne 0 ]; then
             exit 1
     fi
}

function boxUp {
	ping -c 1 -s 10 $1 1>/dev/null 2>&1
	if [ $? != 0 ]; then
		bar="danger"
	else
		bar="success"
	fi
	echo $bar
}

function checkDisk {
	if (( $1 >= 85 )); then
		bar="warning"
	else 
		bar="success"
	fi
	echo "$bar"
}

function checkCPU {
	if (( $1 >= 80 )); then 
		bar="warning"
	else
		bar="success"
	fi
	echo "$bar"
}

function checkBufferMem {
	if (( $1 >= 80 ));then
                bar="warning"
        else
                bar="success"
        fi
        echo "$bar"
}

function checkURL {
	 Console=`wget -O /dev/null "$1" 2>&1 | grep -F HTTP | cut -d ' ' -f 6`
        if [ $Console -eq 200 ]; then
		echo "success"
        else
		echo "danger"
        fi
}

# Format HTML ------------------------------------------------------------
header=$(cat <<EOF
<header>
         <div class="navbar navbar-default navbar-fixed-top">
	<div class="container">
        <div class="navbar-header">
        	<a href="http://dvofsbldap002uk.dev.global.organization.org/" class="navbar-brand"><b>PULSE</b></a>
          	<button class="navbar-toggle" type="button" data-toggle="collapse" data-target="#navbar-main">
           	<span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          	</button>
 		</div>	
 		<div class="navbar-collapse collapse" id="navbar-main">
 		<ul class="nav navbar-nav navbar-right">
		</ul>
		</div>
	<div>
</div>			
</header>
EOF
)

table_start=$(cat <<EOF
<div class="container">
<table class="table table-striped table-hover ">
EOF
)

table_header_1=$(cat <<EOF
<thead>
<tr>
<th>Servers</th>
<th>Service</th>
</tr>
</thead>
<tbody>
EOF
)

table_header_2=$(cat <<EOF
<thread>
<tr>
<th>Servers</th>
<th>System</th>
<th>%</th>
</tr>
</thead>
<tbody>
EOF
) 

footer=$(cat <<EOF
<footer class="text-center">
       	       <h1>Labels</h1>
               <div class="bs-component" style="margin-bottom: 10px;">
               <span class="label label-info">Info</span>
               <span class="label label-success">Active</span>
               <span class="label label-danger">Inactive</span>
               <span class="label label-warning">System Warning</span>
               <span class="label label-default">Pulse-app inactive</span>
               </div>
      		<font size="2">Copyright &copy; organization Technology<br></font>
      		<font size="1">Build version: 1.5</font>
</footer>
EOF
)

# Code ------------------------------------------------------------
echo "<html lang=\"en\">" > $report
echo "<head>" >> $report
echo "<link href=\"data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAA////AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARAAAAAAAAABEAAAAAAAAAEQAAAAAAAAARAAAAAAAAABEREAAAAAAAERERAAAAAAARAAEQAAAAABEAARAAAAAAEQABEAAAAAAREREAAAAAABEREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\" rel=\"icon\" type=\"image/x-icon\" >" >> $report
echo "<meta charset=\"utf-8\"><title>Pulse-Environment</title><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" /><link rel=\"stylesheet\" href=\"${f_css_1}\" media=\"screen\"><link rel=\"stylesheet\" href=\"${f_css_2}\"></head>" >> $report
echo "<body id=\"page-top\" class=\"index\">" >> $report
echo "$header" >> $report
echo "$table_start" >> $report
echo "$table_header_1" >> $report 
while read -u10 line
do
  serverType=`echo $line | awk -F"${FSEP}" '{print $1}'`
  serverList=`echo $line | awk -F"${FSEP}" '{print $2}'`
  serviceType=`echo $line | awk -F"${FSEP}" '{print $3}'`
	
	for server in `echo $serverList`
	do
	    boxBar=`boxUp $server`
    	    if [[ $boxBar == "success" ]];then
		case "$serviceType" in
		"URL" ) url=`echo $line | awk -F"${FSEP}" '{print $4}'`
			servicebar=`checkURL $url`	
			echo "<tr>" >> $report
			echo "<td class=\"info\" rowspan=\"1\"><b>$serverType</b><br>$server</td>" >> $report
			echo "<td class=\"$servicebar\">Console<br>$url</td>" >> $report
			echo "</tr>" >> $report
			;;
		"PROCESS" ) agentList=`echo $line | awk -F"${FSEP}" '{print $4}'`
			   numberOfAgents=`echo $agentList | wc -w` 
			   echo "<tr><td class=\"info\" rowspan=\"$numberOfAgents\"><b>$serverType</b><br>$server</td>" >> $report
			   for agent in `echo $agentList`
			   do
				agentRunning=`ssh $USER@$server ps -ef | grep -v grep | grep -c $agent`
				if [[ "$agentRunning" = "0" ]]; then
					agentBar="danger"
				else 
					agentBar="success"	
				fi
				echo "<td class=\"$agentBar\">$agent</td></tr>" >> $report	
			   done 
			;;		
		esac
	    else
		echo "<tr>" >> $report
        	echo "<td class=\"info\" rowspan=\"1\"><b>$serverType</b><br>$server</td>" >> $report
        	echo "<td class=\"$boxBar\">Machine Down</td></tr>" >> $report
	    fi
	done
done 10< $f_config
echo "</tbody></table><br>" >> $report 

echo "$table_start" >> $report
echo "$table_header_2" >> $report
while read -u10 line
do
  serverType=`echo $line | awk -F"${FSEP}" '{print $1}'`
  serverList=`echo $line | awk -F"${FSEP}" '{print $2}'`

  for server in `echo $serverList`
  do
    boxBar=`boxUp $server`
    if [[ $boxBar == "success" ]];then 
	diskSpace=`ssh $USER@$server df -h /u01 | grep -o [0-9][0-9]*% | sed 's|%||g' | cut -d "." -f1`
	statusDisk=`checkDisk $diskSpace`
	f_Diskbar="${dirNexus}/${statusDisk}.png"
	diskWidth=`awk -vp=$diskSpace -vq=$const_width 'BEGIN{printf "%.2f" ,p * q+10}'`
	
	cpuUtil=`ssh $USER@$server ps aux | awk {'sum+=$3;print sum'} | tail -n 1 | cut -d "." -f1`
	statusCpu=`checkCPU $cpuUtil`
	f_Cpubar="${dirNexus}/${statusCpu}.png"
	cpuWidth=`awk -vp=$cpuUtil -vq=$const_width 'BEGIN{printf "%.2f" ,p * q+10}'`
	
	freeMem=`ssh $USER@$server free | grep Mem | awk '{ printf("%.0f", $4/$2 * 100.0) }'`
	statusMem=`checkBufferMem $freeMem`
	f_Membar="${dirNexus}/${statusMem}.png"
	freeWidth=`awk -vp=$freeMem -vq=$const_width 'BEGIN{printf "%.2f" ,p * q+10}'`

	echo "<tr>" >> $report
		echo "<td class=\"info\" rowspan=\"3\"><b>$serverType</b><br>$server</td>" >> $report
        	echo "<td class=\"$statusDisk\">Diskspace</td>" >> $report
		echo "<td><img src=\"$f_Diskbar\" alt=\"\" width=\"$diskWidth\"  height=\"16\"/>$diskSpace%</td>" >> $report
	echo "</tr>" >> $report
        	echo "<td class=\"$statusCpu\">Cpu Utilisation</td><td class=\"value first\"><img src=\"$f_Cpubar\" alt=\"\" width=\"$cpuWidth\" height=\"16\"/>$cpuUtil%</td>" >> $report
	echo "</tr>" >> $report
        	echo "<td class=\"$statusMem\">Memory Usage</td><td class=\"value first\"><img src=\"$f_Membar\" alt=\"\" width=\"$freeWidth\" height=\"16\"/>$freeMem%</td>" >> $report
	echo "</tr>" >> $report
    else 
	echo "<tr>" >> $report
	echo "<td class=\"info\" rowspan=\"1\"><b>$serverType</b><br>$server</td>" >> $report
	echo "<td class=\"$boxBar\">Machine Down</td></tr>" >> $report
    fi
  done
done 10< $f_config
echo "</tbody></table><br>" >> $report 
echo "<p align=\"centre\">Generated at: <b>` date +\"%r\"`</b>" >> $report
echo "$footer" >> $report
echo "</div></body></html>" >> $report

##### Upload To Nexus the latest data ####
mkdir -p `dirname ${f_status_htm}`
cp $report ${f_status_htm}

cntWarn=`cat $report | grep -v span | grep -c "warning"`
cntDown=`cat $report | grep -v span | grep -c "danger"`

        if [[ ${cntDown} != 0 ]]; then
                echo "danger" > ${f_status}
        elif [[ ${cntWarn} != 0 ]];then
                echo "warning" > ${f_status} 
        else
                echo "success" > ${f_status} 
        fi


datFile=`basename ${f_status}`
htmlFile=`basename ${f_status_htm}`
uploadNexus ${f_status} ${NEXUS_UPLOAD_URL}/${datFile}
uploadNexus ${f_status_htm} ${NEXUS_UPLOAD_URL}/${htmlFile}



