window.onload = function () {

    var seen = new Map()
    var charts = new Map()
    var conn;
    var log = document.getElementById("measurements");

    //
    //
    //
    var fade = function (node) {
        var level = 1;

        var step = function () {
            var hex = level.toString(16);
            node.style.backgroundColor = '#FFFF' + hex + hex;
            if (level < 15) {
                level += 1;
                setTimeout(step, 100);
            }
        }
        setTimeout(step, 100);
    }

    //
    //
    //
    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        //log.append(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    // https://jsfiddle.net/hibbard_eu/C2heg/
    var $divs = $("div.measurement");

    var numericallyOrderedDivs = $divs.sort(function (a, b) {
        return $(a).find("p").text() > $(b).find("p").text();
    });

    //
    //
    //
    //function getMeasurementDIV(obj) {
    const getMeasurementDIV = (obj) => {

        var div,
            alias = obj.alias + "@" + obj.host;

        if (!seen.has(alias)) {

            seen.set(alias, obj);

            const $div = $(createMeasurementLayout(obj));

            div = $div[0]

            appendLog(div);

            charts.set(alias, chartService("div_g_" + alias));

        } else {
            div = document.getElementById(alias)
        }

        return div

    }

    //
    //
    //
    const processMessage = (obj) => {

        var div = getMeasurementDIV(obj)

        div.childNodes[3].innerText = " " + obj.temp + "°C";

        fade(div.childNodes[3])

        charts.get(obj.alias + "@" + obj.host).pushValue(obj)

        //$("#measurements").html(numericallyOrderedDivs);

    }

    //
    //
    //
    const createMeasurementLayout = (obj) => {

        var id = obj.alias + "@" + obj.host,
            alias = obj.alias.replace(/[_]/g, ' ');

        alias = alias + " @ " + obj.host

        return `
            <div id="${id}" class="measurement">
                <p id= "alias" style="color: blue">${alias}</p>
                <div style="border-color: black; border-style: solid" >
                    <p id="values"></p>
                </div>
                <div id="div_g_${id}" style="height:100px;" ></div>
            </div>
          `;
    };

    ////////////////////////////////////////////////////////////////////////////////

    //
    //
    // @see:     // http://code.shutterstock.com/rickshaw/examples/extensions.html
    function chartService(divi) {

        return (function (div) {

            var data = [],
                dygraph = null,
                pushStrategy = null,
                options = {}


            // options = {
            // labelsDivStyles: {
            //     border: '1px solid black'
            // },
            // title: 'Chart Title',
            //     xlabel: 'Date',
            //     ylabel: 'Temperature (F)'
            // }

            //options = {
            // drawPoints: "true",
            // showRoller: "true",
            // valueRange: [0.0, 40.0],
            //labels: ['Time', 'Temperature', 'Hum']
            //                    ylabel: 'Temp',
            //                    y2label: 'Secondary y-axis',
            //  }
            // options = {
            //     labels: ['Date', 'Temp', 'Humidity'],
            //     ylabel: 'Temp',
            //     y2label: 'Hum',
            //     series: {
            //         'Humidity': {
            //             axis: 'y2'
            //         }
            //     },
            //     axes: {
            //         y: {
            //             // set axis-related properties here
            //             drawGrid: true,
            //             independentTicks: false
            //         },
            //         y2: {
            //             // set axis-related properties here
            //             labelsKMB: true,
            //             drawGrid: true,
            //             independentTicks: true
            //         }
            //     }
            // }

            const phenomenonTime = (obj) => (
                new Date(obj.phenomenontime * 1000)
            )

            const setPushStrategy = (obj) => {

                var time = phenomenonTime(obj);

                if (obj.humidity < 0) {

                    data.push([time, 0]);

                    pushStrategy = hasNoHumidity

                } else {

                    data.push([time, 0, 0]);

                    pushStrategy = hasHumidity

                }

            }

            const hasHumidity = (obj) => {
                data.push([phenomenonTime(obj), obj.temp, obj.humidity]);
            }

            const hasNoHumidity = (obj) => {
                data.push([phenomenonTime(obj), obj.temp]);
            }

            return {

                pushValue: function (obj) {

                    if (dygraph === null) {
                        setPushStrategy(obj)
                        dygraph = new Dygraph(div, data, options);
                    }

                    pushStrategy(obj)

                    dygraph.updateOptions({
                        'file': data
                    });
                }
            };
        }(divi))
    }

    // ! ENVELOP einführen !

    if (window["WebSocket"]) {

        conn = new WebSocket("ws://" + document.location.host + "/ws");

        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };

        conn.onmessage = function (evt) {

            var message = evt.data.split('\n');

            for (var i = 0; i < message.length; i++) {
                processMessage(JSON.parse(message[i]))
            }

        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
