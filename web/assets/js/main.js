function gd(year, month, day) {
    return new Date(year, month - 1, day).getTime();
}
var dataStore = {
    "networkSend": {},
    "networkRecv": {},
    "networkSendPacket": {},
    "networkRecvPacket": {},
    "lastNetworkValue": []
}

var initStatus = {
    networkInitStatus: 0,
    diskInitStatus: 0,
    memoryInitStatus: 0,
}

// 管理所有已初始化的组件对象
var visComponentsManager = {}
var MaxPoint = 50
var DNSRefreshTime = 5
var networkPlot = null
var queryInterval = 5000
var now = new Date().getTime();

var gauge_settings = {
    lines: 12,
    angle: 0,
    lineWidth: 0.4,
    pointer: {
        length: 0.75,
        strokeWidth: 0.042,
        color: '#1D212A'
    },
    limitMax: 'false',
    colorStart: '#1ABC9C',
    colorStop: '#1ABC9C',
    strokeColor: '#F0F3F3',
    generateGradient: true,
    highDpiSupport: true, // High resolution support
};
var networkOptions = {
    series: {
        lines: {
            show: true,
            lineWidth: .8,
            fill: true
        }
    },
    xaxis: {
        mode: "time",
        tickSize: [10, "second"],
        tickFormatter: function(v, axis) {
            var date = new Date(v);
            if (date.getSeconds() % 60 == 0) {
                var hours = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
                var minutes = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
                var seconds = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();

                return hours + ":" + minutes + ":" + seconds;
            } else {
                return "";
            }
        },
        axisLabelUseCanvas: true,
        axisLabelFontSizePixels: 10,
        axisLabelFontFamily: 'Verdana, Arial',
        axisLabelPadding: 2
    },
    yaxis: {
        axisLabel: "实时网络流量",
        axisLabelUseCanvas: true,
        axisLabelFontSizePixels: 12,
        axisLabelFontFamily: 'Verdana, Arial',
        axisLabelPadding: 6,
        min: 0,
    },
    grid: {
        hoverable: true,
        borderWidth: 2,
        backgroundColor: {
            colors: ["#EDF5FF", "#ffffff"]
        }
    },
    tooltip: {
        show: true,
        content: "接口:%s %y"
    },
    legend: {
        labelBoxBorderColor: "white",
        position: "nw",
        backgroundColor: "grey",
        backgroundOpacity: 0.4

    }
};


var dnsQueryOptions = {
    series: {
        stack: true,
        lines: {
            show: true,
            barWidth: .8,
            fill: true
        }
    },
    grid: {
        hoverable: true,
        // backgroundColor: { colors: ["#EDF5FF", "#ffffff"] }
    },
    tooltip: {
        show: true,
        content: '<h4>%s</h4>查询量:%y'
    },
    xaxis: {
        mode: "time",
        tickSize: [5, "second"],
        tickFormatter: function(v, axis) {
            var date = new Date(v);
            if (date.getSeconds() % 60 == 0) {
                var hours = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
                var minutes = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
                var seconds = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();

                return hours + ":" + minutes + ":" + seconds;
            } else {
                return "";
            }
        },
        axisLabelUseCanvas: true,
        axisLabelFontSizePixels: 10,
        axisLabelFontFamily: 'Verdana, Arial',
        axisLabelPadding: 2
    },
    yaxis: {
        axisLabel: "实时DNS查询流量",
        axisLabelUseCanvas: true,
        axisLabelFontSizePixels: 12,
        axisLabelFontFamily: 'Verdana, Arial',
        axisLabelPadding: 6,
        min: 0,
    },
    legend: {
        labelBoxBorderColor: "white",
        position: "nw",
        backgroundColor: "grey",
        backgroundOpacity: 0.4
    }
};

var optionBytesSend = {
    bps: true,
    axisLabel: "实时发送数据字节速率"
}
var optionBytesRecv = {
    bps: true,
    axisLabel: "实时接收数据字节速率"
}
var optionPacketsSend = {
    bps: false,
    axisLabel: "实时发送数据包速率"
}
var optionPacketsRecv = {
    bps: false,
    axisLabel: "实时接收数据包速率"
}

var pie_options = {
    series: {
        pie: {
            show: true,
            innerRadius: 0.5,
            radius: 1
        }
    },
    grid: {
        hoverable: true
    },
    tooltip: true,
    tooltipOpts: {
        cssClass: "flotTip",
        content: "%s: %p.0%",
        defaultTheme: false
    }
};

function tickNetworkFormatter(bps) {
    return function(v) {
        if (bps === true) {
            v = v * 8
        }
        v = v / (queryInterval / 1000)
        if (v > (1024 * 1024 * 1024)) {
            return (v / (1024 * 1024 * 1024)).toFixed(0) + (bps ? "Gbps" : "Gpps");
        } else if (v > (1024 * 1024)) {
            return (v / (1024 * 1024)).toFixed(0) + (bps ? "Mbps" : "Mpps");
        } else if (v > (1024)) {
            return (v / (1024)).toFixed(0) + (bps ? "Kbps" : "Kpps");
        } else {
            return v + (bps ? "bps" : "pps");
        }
    }
}

function formatMemoryUsage(v) {
    if (v > (1024 * 1024 * 1024)) {
        return (v / (1024 * 1024 * 1024)).toFixed(1) + "GB";
    } else if (v > (1024 * 1024)) {
        return (v / (1024 * 1024)).toFixed(1) + "MB";
    } else if (v > (1024)) {
        return (v / (1024)).toFixed(1) + "KB";
    } else {
        return v;
    }
}

function init_network_chart(data, selector, storeName, trafficName, option) {
    if (typeof data.network_io === 'undefined' || data.network_io.length === 0) {
        console.log("network info not ready")
        return
    }
    var networkCardsNameList = data.network_io.map(function(item) {
        // 初始化数组
        dataStore[storeName][item.name] = []
        for (var i = 0; i < MaxPoint; i++) {
            dataStore[storeName][item.name].push([now - (MaxPoint - i) * queryInterval, null])
        }
        dataStore["lastNetworkValue"][item.name] = _.clone(item, true)
        return item.name;
    })

    var dataSet = []
    for (var i = 0; i < networkCardsNameList.length; i++) {
        dataSet.push({
            label: networkCardsNameList[i],
            data: dataStore[storeName][networkCardsNameList[i]],
            color: BasicColorSets[i],
        })
    }
    networkOptions.yaxis.axisLabel = option.axisLabel
    networkOptions.yaxis.tickFormatter = tickNetworkFormatter(option.bps)
    $.plot(selector, dataSet, networkOptions);
    initStatus.networkInitStatus += 1
}

function update_network_chart(data, selector, storeName, trafficName, updateValue, option) {
    if (initStatus.networkInitStatus < 4) {
        init_network_chart(data, selector, storeName, trafficName, option)
        return
    }
    var networkCardsNameList = data.network_io.map(function(item) {
        // 初始化数组
        dataStore[storeName][item.name].shift();
        var value = item[trafficName] - dataStore["lastNetworkValue"][item.name][trafficName]
        if (value < 0) {
            value = null
        }
        dataStore[storeName][item.name].push([now, value])
        if (updateValue === true) {
            // 不会再使用上一次的数值
            dataStore["lastNetworkValue"][item.name] = _.clone(item, true)
        }
        return item.name;
    })
    var dataSet = []
    for (var i = 0; i < networkCardsNameList.length; i++) {
        dataSet.push({
            label: networkCardsNameList[i],
            data: dataStore[storeName][networkCardsNameList[i]],
            color: BasicColorSets[i],
        })
    }
    networkOptions.yaxis.axisLabel = option.axisLabel
    networkOptions.yaxis.tickFormatter = tickNetworkFormatter(option.bps)
    $.plot(selector, dataSet, networkOptions);
}

// 格式化磁盘数据到一个object , 每一个key对应一个挂载点，value为flot的 pie option结构
function formatRawDiskDataToOption(diskInfoObject) {
    var result = {}
    for (var i = 0; i < diskInfoObject.length; i++) {
        result[diskInfoObject[i]['path']] = [{
            label: "未用",
            data: parseInt(diskInfoObject[i]['free'])
        }, {
            label: "已使用",
            data: parseInt(diskInfoObject[i]['used'])
        }]
    }
    return result
}


function update_disk_chart(data) {
    if (typeof data.disk_usage === 'undefined' || data.disk_usage.length === 0) {
        console.warn("disk info not ready")
        return
    }
    // format data to dataset
    var dataSets = formatRawDiskDataToOption(data.disk_usage)
    var diskNameList = _.keys(dataSets)
    for (var i = 0; i < diskNameList.length; i++) {
        var divIdName = "disk_plot_01" + "_" + i
        if (initStatus.diskInitStatus < diskNameList.length) {
            $("#disk_plot_01").append('<div class="col-md-3 col-sm-6"><div class="title">' + diskNameList[i] + '</div><div class="disk-item" id=' + divIdName + '></div></div>');
            initStatus.diskInitStatus += 1
        }
        var divIdNameSelect = '#' + divIdName
        $.plot($(divIdNameSelect), dataSets[diskNameList[i]], pie_options);
    }
}

function update_process_chart(data, selectorID, listTitle) {
    if (!data || typeof data[listTitle] === 'undefined' || typeof data[listTitle].length === 0) {
        console.log('unable to update process info')
        return
    }
    $(selectorID).html('')
    var htmlContent = ''
    for (var i = 0; i < data[listTitle].length; i++) {
        var classEvenOrOdd = (i % 2 == 0) ? "even pointer" : "odd pointer";
        htmlContent += '<tr class="' + classEvenOrOdd + '">' +
            '<td class=" ">' + data[listTitle][i]['pid'] + '</td>' +
            '<td class=" ">' + data[listTitle][i]['ppid'] + '</td>' +
            '<td class=" ">' + data[listTitle][i]['command'] + '</td>' +
            '<td class=" ">' + data[listTitle][i]['memory_percent'] + '</td>' +
            '<td class=" ">' + data[listTitle][i]['cpu_percent'] + '</td></tr>';
    }
    $(selectorID).append(htmlContent)
}

function parseStatusData(data) {
    // 更新系统基本信息
    $(".loading").fadeOut()
    baiscStatus(data)
    loadStatus(data)
    now += queryInterval
        // 更新网络接口图示
    update_network_chart(data, $("#network_plot_01"), "networkSend", "bytesSent", false, optionBytesSend);
    update_network_chart(data, $("#network_plot_02"), "networkSendPacket", "packetsSent", false, optionPacketsSend);
    update_network_chart(data, $("#network_plot_03"), "networkRecv", "bytesRecv", false, optionBytesRecv);
    update_network_chart(data, $("#network_plot_04"), "networkRecvPacket", "packetsRecv", true, optionPacketsRecv);

    // 更新磁盘使用图示
    update_disk_chart(data)
        // 更新进程状态
    update_process_chart(data, '#system-process-memory-list', 'processes_memory')
    update_process_chart(data, '#system-process-cpu-list', 'processes_cpu')

}

function baiscStatus(data) {
    if (typeof data.host_info === 'undefined') {
        return
    }
    if (data.host_info.hostname) {
        $("#info-hostname").html(data.host_info.hostname)
    } else {
        $("#info-hostname").html("unknown")
    }
    if (data.host_info.uptime) {
        $("#info-uptime").html(formatTimeSeconds(data.host_info.uptime))
    } else {
        $("#info-uptime").html("unknown")
    }
    if (data.host_info.procs) {
        $("#info-procs").html((data.host_info.procs))
    } else {
        $("#info-procs").html("unknown")
    }
    if (data.host_info.kernelVersion) {
        $("#info-kernel").html(data.host_info.kernelVersion.split("-")[0])
        $("#info-kernel-detail").html(data.host_info.kernelVersion)
    } else {
        $("#info-kernel").html("unknown")
        $("#info-kernel-detail").html("unknown")
    }
    if (data.host_info.platform) {
        $("#info-platform").html(data.host_info.platform)
    } else {
        $("#info-platform").html("unknown")
    }
    if (data.host_info.platformVersion) {
        $("#info-platform-version").html(data.host_info.platformVersion)
    } else {
        $("#info-platform-version").html("unknown")
    }
    // update memory info 
    if (typeof(Gauge) === 'undefined' || typeof data["virtual_memory"] === 'undefined') {
        return;
    }
    if ($('#info-system-memory').length) {
        var chart_gauge_memory = document.getElementById('info-system-memory');
        var gauge_memory
        var used = parseFloat(data["virtual_memory"]["used"])
        var total = parseFloat(data["virtual_memory"]["total"])
        if (initStatus.memoryInitStatus === 0) {
            gauge_memory = new Gauge(chart_gauge_memory).setOptions(gauge_settings);
            gauge_memory.maxValue = total
            gauge_memory.animationSpeed = 32;
            visComponentsManager['memoryMeter'] = gauge_memory
            initStatus.memoryInitStatus = 1
        }

        visComponentsManager['memoryMeter'].set(used)
        $("#hover-for-memory").attr('title', formatMemoryUsage(used) + '/' + formatMemoryUsage(total))
    }
    return
}

function loadStatus(data) {
    if (typeof data.load_avg === 'undefined') {
        return
    }
    if (data.load_avg.load1) {
        $("#info-load-1min").html(data.load_avg.load1 + "%")
    } else {
        $("#info-load-1min").html("")
    }
    if (data.load_avg.load5) {
        $("#info-load-5min").html(data.load_avg.load5 + "%")
    } else {
        $("#info-load-5min").html("")
    }
    if (data.load_avg.load15) {
        $("#info-load-15min").html(data.load_avg.load15 + "%")
    } else {
        $("#info-load-15min").html("")
    }
}

function float5SecondTimeStamp(current) {
    return new Date(current.getTime() - current.getTime() % 5000)
}

function parseInitDNSStatus(data) {
    var originalDataMap = {}
    var now = float5SecondTimeStamp(new Date()).getTime()
    if (typeof data === 'undefined' || data.length === 0) {
        console.warn("dns query history data not ready")
        return
    }

    data.map(function(item) {
        var temp = item.timestamp.split(' ')
        var timestamp = float5SecondTimeStamp(new Date(temp[0] + ' ' + temp[1])).getTime()
        if (dataStore['lastDNSUpdateTime']) {
            if (timestamp > dataStore['lastDNSUpdateTime']) {
                dataStore['lastDNSUpdateTime'] = timestamp
            }
        } else {
            // 未初始化
            dataStore['lastDNSUpdateTime'] = timestamp
        }
        if (typeof item['type_stats'] !== 'string') {
            return
        }
        var typeinfo = JSON.parse(item['type_stats']) || []

        // 初始化默认数组
        for (var i = 0; i < typeinfo.length; i++) {
            if (typeof originalDataMap[typeinfo[i]['type']] === 'undefined') {
                // 初始化一组类型数据并设置为空值
                originalDataMap[typeinfo[i]['type']] = []
                for (var j = 0; j < MaxPoint; j++) {
                    var tempDate = now - (MaxPoint - j + 1) * 5000
                    var value = (timestamp === tempDate) ? Math.ceil(typeinfo[i]['sum'] / 5) : null
                    originalDataMap[typeinfo[i]['type']].push([tempDate, value])
                }
            } else {
                var typeDataList = originalDataMap[typeinfo[i]['type']]
                for (var j = 0; j < typeDataList.length; j++) {
                    if (timestamp === typeDataList[j][0]) {
                        typeDataList[j][1] = Math.ceil(typeinfo[i]['sum'] / DNSRefreshTime)
                    }
                }
            }
        }
    })
    var serialData = []
    _.keys(originalDataMap).map(function(name) {
        serialData.push({
            data: originalDataMap[name],
            label: name
        })
    });
    dataStore['dnsSerialData'] = serialData
    $.plot("#dns-realtime-query", serialData, dnsQueryOptions);
    return serialData
}

function UpdateDNSStatus(data) {
    if (typeof data !== 'object' || typeof data['type_stats'] === 'undefined') {
        console.warn("dns query data not ready")
        return
    }
    var serialData = dataStore.dnsSerialData
    var temp = data['timestamp'].split(' ')
    var timestamp = float5SecondTimeStamp(new Date(temp[0] + ' ' + temp[1])).getTime()
    if (timestamp <= dataStore.lastDNSUpdateTime) {
        return
    }

    var typeDataList = JSON.parse(data['type_stats']) || []
    if (typeDataList.length === 0) {
        // 暂时没有任何数据获取到，但是仍旧需要更新图表
        for (var j = 0; j < serialData.length; j++) {
            if (serialData[j]['data'].length === MaxPoint) {
                serialData[j]["data"].shift()
            }
            if (dataStore.lastDNSUpdateTime < timestamp) {
                dataStore.lastDNSUpdateTime = timestamp
            }
            serialData[j]["data"].push([timestamp, 0])
        }

    } else {
        for (var i = 0; i < typeDataList.length; i++) {
            var found = false
            for (var j = 0; j < serialData.length; j++) {
                if (serialData[j]['label'] === typeDataList[i]['type']) {
                    if (serialData[j]['data'].length === MaxPoint) {
                        serialData[j]["data"].shift()
                    }
                    if (dataStore.lastDNSUpdateTime < timestamp) {
                        dataStore.lastDNSUpdateTime = timestamp
                    }
                    serialData[j]["data"].push([timestamp, Math.ceil(typeDataList[i]['sum'] / DNSRefreshTime)])
                    found = true
                }
            }
            if (!found) {
                var newSerialData = {
                    label: typeDataList[i]['type'],
                    data: [
                        [now, Math.ceil(typeDataList[i]['sum'] / DNSRefreshTime)]
                    ]
                }
                serialData.push(newSerialData)
            }
        }
    }
    $.plot("#dns-realtime-query", serialData, dnsQueryOptions);
}

function UpdateTableList(dataList, selectorId, keyName, valueName) {
    var selector = $(selectorId);
    selector.html('')
    var htmlContent = ''
    for (var i = 0; i < dataList.length; i++) {
        htmlContent += '<tr>' +
            '<td>' + (i + 1) + '</td>' +
            '<td>' + dataList[i][keyName] + '</td>' +
            '<td>' + Math.ceil(dataList[i][valueName] / DNSRefreshTime) + '</td></tr>';
    }
    selector.append(htmlContent)
}
var topPieOption = {
    series: {
        pie: {
            show: true,
            radius: 1,
            label: {
                show: false,
            },
            innerRadius: 0.5,
        }
    },
    legend: {
        show: false
    },
    grid: {
        hoverable: true
    },
    tooltip: true,
    tooltipOpts: {
        cssClass: "flotTip",
        content: "%s: %p.0%",
        defaultTheme: false
    }
}

function UpdatePieComponent(dataSet, selectorId, keyName, valueName) {
    var data = []
    var total = 0
    for (var i = 0; i < dataSet.length; i++) {
        data.push({
            label: dataSet[i][keyName],
            data: parseInt(dataSet[i][valueName]),
        })
        total += parseInt(dataSet[i][valueName]);
    }
    var unknown = dataStore.currentTotal - total
    if (unknown > 0) {
        data.push({
            label: "其他",
            data: unknown
        })
    }
    $.plot($(selectorId), data, topPieOption);
}


function UpdateTopStatus(data) {
    if (typeof data !== 'object') {
        return
    }
    var type_stats = JSON.parse(data['type_stats'])
    if (typeof data['type_stats'] === 'undefined' || type_stats === null || type_stats.length === 0) {
        $('#type-top-cover').show();
    } else {
        if ($('#type-top-cover').is(':visible')) {
            $('#type-top-cover').hide();
        }
        var total = 0
        for (var i = 0; i < type_stats.length; i++) {
            total += type_stats[i]['sum']
        }
        dataStore.currentTotal = total;
        // 更新数据表
        UpdateTableList(type_stats, '#type-top-table', 'type', 'sum')
            // 更新pie饼图
        UpdatePieComponent(type_stats, '#type-top-doughnut', 'type', 'sum')
    }


    var ip_stats = JSON.parse(data['ip_stats'])
    if (typeof data['ip_stats'] === 'undefined' || ip_stats === null || ip_stats.length === 0) {
        $('#ip-top-cover').show();
    } else {
        if ($('#ip-top-cover').is(':visible')) {
            $('#ip-top-cover').hide();
        }
        UpdateTableList(ip_stats, '#ip-top-table', 'ip', 'sum')

        // 更新pie饼图
        UpdatePieComponent(ip_stats, '#ip-top-doughnut', 'ip', 'sum')
    }



    var domain_stats = JSON.parse(data['domain_stats'])
    if (typeof data['domain_stats'] === 'undefined' || domain_stats === null || domain_stats.length === 0) {
        $('#domain-top-cover').show();

    } else {
        if ($('#domain-top-cover').is(':visible')) {
            $('#domain-top-cover').hide();
        }
        UpdateTableList(domain_stats, '#domain-top-table', 'domain', 'sum')

        // 更新pie饼图
        UpdatePieComponent(domain_stats, '#domain-top-doughnut', 'domain', 'sum')

    }

    var sub_domain_stats = JSON.parse(data['sub_domain_stats'])
    if (typeof data['sub_domain_stats'] === 'undefined' || sub_domain_stats === null || sub_domain_stats.length === 0) {
        $('#tld-domain-top-cover').show();

    } else {
        if ($('#tld-domain-top-cover').is(':visible')) {
            $('#tld-domain-top-cover').hide();
        }
        UpdateTableList(sub_domain_stats, '#tld-domain-top-table', 'domain', 'sum')
            // 更新pie饼图
        UpdatePieComponent(sub_domain_stats, '#tld-domain-top-doughnut', 'domain', 'sum')

    }

}

function formatTimeSeconds(seconds) {
    if (seconds > 86400) {
        return (seconds / 86400).toFixed(2).toString() + " d"
    }
    if (seconds > 3600) {
        return (seconds / 3600).toFixed(2).toString() + "h"
    }

    if (seconds > 60) {
        return (seconds / 60).toFixed(2).toString() + "m"
    }
    return seconds.toString() + "s"
}
$(document).ready(function() { 
    init_sidebar()
    $.ajax({
        dataType: "json",
        url: "/dns_init_status",
        success: function(data) {
            if (typeof data !== "object" || typeof data['data'] !== 'object' || data['error'] !== '') {
                new PNotify({
                    title: '错误提示',
                    styling: 'bootstrap3',
                    type: "error",
                    text: data['error'] || '服务暂时无法访问，请稍后再试'
                });
                return
            }
            parseInitDNSStatus(data['data'])
        },
        error: function(error) {
            new PNotify({
                title: '错误提示',
                styling: 'bootstrap3',
                type: "error",
                text: error['error'] || '服务暂时无法访问，请稍后再试'
            });
        }
    }).fail(function(jqXHR, textStatus, errorThrown) {
        new PNotify({
            title: '错误提示',
            styling: 'bootstrap3',
            type: "error",
            text: '网络服务中断，请检查网络连接是否正常'
        });
    });
    setInterval(function() {
        $.ajax({
            dataType: "json",
            url: "/status",
            success: function(data) {
                if (typeof data !== "object" || typeof data['data'] !== 'object' || data['error'] !== '') {
                    new PNotify({
                        title: '数据获取失败',
                        styling: 'bootstrap3',
                        type: "error",
                        text: data['error'] || '服务暂时无法访问，请稍后再试'
                    });
                    return
                }
                parseStatusData(data['data'])
            },
            error: function(error) {
                new PNotify({
                    title: '错误提示',
                    styling: 'bootstrap3',
                    type: "error",
                    text: error['error'] || '服务暂时无法访问，请稍后再试'
                });
            }
        }).fail(function() {
            new PNotify({
                title: '错误提示',
                styling: 'bootstrap3',
                type: "error",
                text: '网络服务中断，请检查网络连接是否正常'
            });
        });
        $.ajax({
            dataType: "json",
            url: "/dns_lastest_status",
            success: function(data) {
                if (typeof data !== "object" || typeof data['data'] !== 'object' || data['error'] !== '') {
                    new PNotify({
                        title: '错误提示',
                        styling: 'bootstrap3',
                        type: "error",
                        text: data["error"] || '服务暂时无法访问，请稍后再试'
                    });
                    return
                }
                UpdateDNSStatus(data['data'])
                UpdateTopStatus(data['data'])
            },
            error: function(error) {
                new PNotify({
                    title: '错误提示',
                    styling: 'bootstrap3',
                    type: "error",
                    text: error['error'] || '服务暂时无法访问，请稍后再试'
                });
            }

        }).fail(function(jqXHR, textStatus, errorThrown) {
            new PNotify({
                title: '错误提示',
                styling: 'bootstrap3',
                type: "error",
                text: '网络服务中断，请检查网络连接是否正常'
            });
        });
    }, queryInterval)
});