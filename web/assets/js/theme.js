/**
 * Resize function without multiple trigger
 * 
 * Usage:
 * $(window).smartresize(function(){  
 *     // code here
 * });
 */
(function($, sr) {
    // debouncing function from John Hann
    // http://unscriptable.com/index.php/2009/03/20/debouncing-javascript-methods/
    var debounce = function(func, threshold, execAsap) {
        var timeout;

        return function debounced() {
            var obj = this,
                args = arguments;

            function delayed() {
                if (!execAsap)
                    func.apply(obj, args);
                timeout = null;
            }

            if (timeout)
                clearTimeout(timeout);
            else if (execAsap)
                func.apply(obj, args);

            timeout = setTimeout(delayed, threshold || 100);
        };
    };

    // smartresize 
    jQuery.fn[sr] = function(fn) { return fn ? this.bind('resize', debounce(fn)) : this.trigger(sr); };

})(jQuery, 'smartresize');
/**
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

var CURRENT_URL = window.location.href.split('#')[0].split('?')[0],
    $BODY = $('body'),
    $MENU_TOGGLE = $('#menu_toggle'),
    $SIDEBAR_MENU = $('#sidebar-menu'),
    $SIDEBAR_FOOTER = $('.sidebar-footer'),
    $LEFT_COL = $('.left_col'),
    $RIGHT_COL = $('.right_col'),
    $NAV_MENU = $('.nav_menu'),
    $FOOTER = $('footer');

// 产生随机颜色
function getRandomColor() {
    var letters = '0123456789ABCDEF';
    var color = '#';
    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}
var BasicColorSets = []
for (var i = 0; i < 100; i++) {
    BasicColorSets.push(getRandomColor())
}
// Sidebar
function init_sidebar() {
    // TODO: This is some kind of easy fix, maybe we can improve this
    var setContentHeight = function() {
        // reset height
        $RIGHT_COL.css('min-height', $(window).height());

        var bodyHeight = $BODY.outerHeight(),
            footerHeight = $BODY.hasClass('footer_fixed') ? -10 : $FOOTER.height(),
            leftColHeight = $LEFT_COL.eq(1).height() + $SIDEBAR_FOOTER.height(),
            contentHeight = bodyHeight < leftColHeight ? leftColHeight : bodyHeight;

        // normalize content
        contentHeight -= $NAV_MENU.height() + footerHeight;

        $RIGHT_COL.css('min-height', contentHeight);
    };

    $SIDEBAR_MENU.find('a').on('click', function(ev) {
        console.log('clicked - sidebar_menu');
        var $li = $(this).parent();

        if ($li.is('.active')) {
            $li.removeClass('active active-sm');
            $('ul:first', $li).slideUp(function() {
                setContentHeight();
            });
        } else {
            // prevent closing menu if we are on child menu
            if (!$li.parent().is('.child_menu')) {
                $SIDEBAR_MENU.find('li').removeClass('active active-sm');
                $SIDEBAR_MENU.find('li ul').slideUp();
            } else {
                if ($BODY.is(".nav-sm")) {
                    $SIDEBAR_MENU.find("li").removeClass("active active-sm");
                    $SIDEBAR_MENU.find("li ul").slideUp();
                }
            }
            $li.addClass('active');

            $('ul:first', $li).slideDown(function() {
                setContentHeight();
            });
        }
    });

    // toggle small or large menu 
    $MENU_TOGGLE.on('click', function() {
        console.log('clicked - menu toggle');

        if ($BODY.hasClass('nav-md')) {
            $SIDEBAR_MENU.find('li.active ul').hide();
            $SIDEBAR_MENU.find('li.active').addClass('active-sm').removeClass('active');
        } else {
            $SIDEBAR_MENU.find('li.active-sm ul').show();
            $SIDEBAR_MENU.find('li.active-sm').addClass('active').removeClass('active-sm');
        }

        $BODY.toggleClass('nav-md nav-sm');

        setContentHeight();

        $('.dataTable').each(function() { $(this).dataTable().fnDraw(); });
    });

    // check active menu
    $SIDEBAR_MENU.find('a[href="' + CURRENT_URL + '"]').parent('li').addClass('current-page');

    $SIDEBAR_MENU.find('a').filter(function() {
        return this.href == CURRENT_URL;
    }).parent('li').addClass('current-page').parents('ul').slideDown(function() {
        setContentHeight();
    }).parent().addClass('active');

    // recompute content when resizing
    $(window).smartresize(function() {
        setContentHeight();
    });

    setContentHeight();

    // fixed sidebar
    if ($.fn.mCustomScrollbar) {
        $('.menu_fixed').mCustomScrollbar({
            autoHideScrollbar: true,
            theme: 'minimal',
            mouseWheel: { preventDefault: true }
        });
    }
};
var randNum = function() {
    return (Math.floor(Math.random() * (1 + 40 - 20))) + 20;
};

// Panel toolbox
$(document).ready(function() {
    $('.collapse-link').on('click', function() {
        var $BOX_PANEL = $(this).closest('.x_panel'),
            $ICON = $(this).find('i'),
            $BOX_CONTENT = $BOX_PANEL.find('.x_content');

        // fix for some div with hardcoded fix class
        if ($BOX_PANEL.attr('style')) {
            $BOX_CONTENT.slideToggle(200, function() {
                $BOX_PANEL.removeAttr('style');
            });
        } else {
            $BOX_CONTENT.slideToggle(200);
            $BOX_PANEL.css('height', 'auto');
        }

        $ICON.toggleClass('fa-chevron-up fa-chevron-down');
    });

    $('.close-link').click(function() {
        var $BOX_PANEL = $(this).closest('.x_panel');

        $BOX_PANEL.remove();
    });
});
// /Panel toolbox

// Tooltip
$(document).ready(function() {
    $('[data-toggle="tooltip"]').tooltip({
        container: 'body'
    });
});
// /Tooltip

// Progressbar
if ($(".progress .progress-bar")[0]) {
    $('.progress .progress-bar').progressbar();
}
// /Progressbar

// Switchery
$(document).ready(function() {
    if ($(".js-switch")[0]) {
        var elems = Array.prototype.slice.call(document.querySelectorAll('.js-switch'));
        elems.forEach(function(html) {
            var switchery = new Switchery(html, {
                color: '#26B99A'
            });
        });
    }
});
// /Switchery


// iCheck
$(document).ready(function() {
    if ($("input.flat")[0]) {
        $(document).ready(function() {
            $('input.flat').iCheck({
                checkboxClass: 'icheckbox_flat-green',
                radioClass: 'iradio_flat-green'
            });
        });
    }
});
// /iCheck

// Table
$('table input').on('ifChecked', function() {
    checkState = '';
    $(this).parent().parent().parent().addClass('selected');
    countChecked();
});
$('table input').on('ifUnchecked', function() {
    checkState = '';
    $(this).parent().parent().parent().removeClass('selected');
    countChecked();
});

var checkState = '';

$('.bulk_action input').on('ifChecked', function() {
    checkState = '';
    $(this).parent().parent().parent().addClass('selected');
    countChecked();
});
$('.bulk_action input').on('ifUnchecked', function() {
    checkState = '';
    $(this).parent().parent().parent().removeClass('selected');
    countChecked();
});
$('.bulk_action input#check-all').on('ifChecked', function() {
    checkState = 'all';
    countChecked();
});
$('.bulk_action input#check-all').on('ifUnchecked', function() {
    checkState = 'none';
    countChecked();
});

function countChecked() {
    if (checkState === 'all') {
        $(".bulk_action input[name='table_records']").iCheck('check');
    }
    if (checkState === 'none') {
        $(".bulk_action input[name='table_records']").iCheck('uncheck');
    }

    var checkCount = $(".bulk_action input[name='table_records']:checked").length;

    if (checkCount) {
        $('.column-title').hide();
        $('.bulk-actions').show();
        $('.action-cnt').html(checkCount + ' Records Selected');
    } else {
        $('.column-title').show();
        $('.bulk-actions').hide();
    }
}



// Accordion
$(document).ready(function() {
    $(".expand").on("click", function() {
        $(this).next().slideToggle(200);
        $expand = $(this).find(">:first-child");

        if ($expand.text() == "+") {
            $expand.text("-");
        } else {
            $expand.text("+");
        }
    });
});

// NProgress
if (typeof NProgress != 'undefined') {
    $(document).ready(function() {
        NProgress.start();
    });

    $(window).load(function() {
        NProgress.done();
    });
}


//hover and retain popover when on popover content
var originalLeave = $.fn.popover.Constructor.prototype.leave;
$.fn.popover.Constructor.prototype.leave = function(obj) {
    var self = obj instanceof this.constructor ?
        obj : $(obj.currentTarget)[this.type](this.getDelegateOptions()).data('bs.' + this.type);
    var container, timeout;

    originalLeave.call(this, obj);

    if (obj.currentTarget) {
        container = $(obj.currentTarget).siblings('.popover');
        timeout = self.timeout;
        container.one('mouseenter', function() {
            //We entered the actual popover – call off the dogs
            clearTimeout(timeout);
            //Let's monitor popover content instead
            container.one('mouseleave', function() {
                $.fn.popover.Constructor.prototype.leave.call(self, self);
            });
        });
    }
};

$('body').popover({
    selector: '[data-popover]',
    trigger: 'click hover',
    delay: {
        show: 50,
        hide: 400
    }
});


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

var networkInitStatus = 1
var MaxPoint = 50
var networkPlot = null
var queryInterval = 5000
var now = new Date().getTime();

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
        tickSize: [1, "second"],
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
        axisLabelPadding: 6
    },
    legend: {
        labelBoxBorderColor: "white"
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

function init_network_chart(data, selector, storeName, trafficName, option) {
    if (typeof data.network_io === 'undefine' || data.network_io.length === 0) {
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
    networkInitStatus += 1
}

function update_network_chart(data, selector, storeName, trafficName, updateValue, option) {
    if (networkInitStatus < 5) {
        init_network_chart(data, selector, storeName, trafficName, option)
        return
    }
    var networkCardsNameList = data.network_io.map(function(item) {
        // 初始化数组
        dataStore[storeName][item.name].shift();
        if (item.name == 'docker0' && storeName == 'networkSend' && trafficName == 'bytesSent') {
            console.log(item.name, storeName, trafficName, item[trafficName] - dataStore["lastNetworkValue"][item.name][trafficName])
        }
        if (item.name == 'docker0' && storeName == 'networkRecv' && trafficName == 'bytesRecv') {
            console.log(item.name, storeName, trafficName, item[trafficName] - dataStore["lastNetworkValue"][item.name][trafficName])
        }
        dataStore[storeName][item.name].push([now, item[trafficName] - dataStore["lastNetworkValue"][item.name][trafficName]])
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

function parseStatusData(data) {
    baiscStatus(data)
    loadStatus(data)
    now += queryInterval
    update_network_chart(data, $("#network_plot_01"), "networkSend", "bytesSent", false, optionBytesSend);
    update_network_chart(data, $("#network_plot_02"), "networkSendPacket", "packetsSent", false, optionPacketsSend);
    update_network_chart(data, $("#network_plot_03"), "networkRecv", "bytesRecv", false, optionBytesRecv);
    update_network_chart(data, $("#network_plot_04"), "networkRecvPacket", "packetsRecv", true, optionPacketsRecv);
}

function baiscStatus(data) {
    if (typeof data.host_info === 'undefine') {
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
}

function loadStatus(data) {
    if (typeof data.load_avg === 'undefine') {
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

$(document).ready(function() {
    setInterval(function() {
        $.ajax({
            dataType: "json",
            url: "/status",
            success: function(data) {
                parseStatusData(data)
            }
        });
    }, queryInterval)
});

function formatTimeSeconds(seconds) {
    if (seconds > 86400) {
        return (seconds / 86400).toFixed(2).toString() + " days"
    }
    if (seconds > 3600) {
        return (seconds / 3600).toFixed(2).toString() + "hours"
    }

    if (seconds > 60) {
        return (seconds / 60).toFixed(2).toString() + "mins"
    }
    return seconds.toString() + "seconds"
}