// Copyright 2017 The hchart Authors, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

(function (global, undefined) {
    "use strict";

    if (global.hooto_chart) {
        return;
    }

    var hc = (global.hooto_chart = {
        version: "0.1.1",
        basepath: "/hchart/~",
        opts_width: "800px",
        opts_height: "600px",
        instances: {},
    });

    var color = {
        theme: null,
        theme_tb: ["337ab7", "5cb85c", "5bc0de", "f0ad4e", "d9534f", "333333"],
        theme_gc: ["0057e7", "d62d20", "ffa700", "008744", "9c27b0", "333333"],
    };

    function color_pallet(theme) {
        if (!theme) {
            this.theme = color.theme_gc;
        }
        this.randlast = 0;
        this.reset = 1;

        this.rand = function () {
            this.randlast++;
            if (this.randlast > this.theme.length) {
                this.randlast = 1;
                this.reset = this.reset * 0.85;
            }
            var color_i16 = parseInt(this.theme[this.randlast - 1], 16);
            if (this.reset < 1) {
                var r = (color_i16 >> 16) & 255;
                var g = (color_i16 >> 8) & 255;
                var b = color_i16 & 255;
                r = r + parseInt((255 - r) * (1 - this.reset));
                g = g + parseInt((255 - g) * (1 - this.reset));
                b = b + parseInt((255 - b) * (1 - this.reset));
                color_i16 = (1 << 24) + (r << 16) + (g << 8) + b;
            }
            return color_i16;
        };
    }

    color.rgb = function (color_i16, alpha) {
        var r = (color_i16 >> 16) & 255;
        var g = (color_i16 >> 8) & 255;
        var b = color_i16 & 255;

        return "rgba(" + r + "," + g + "," + b + "," + alpha + ")";
    };

    hc.deps_load = function (cb) {
        seajs.use([hc.basepath + "/chartjs/chart.js"], cb);
    };

    hc.JsonRenderElementID = function (elem_id) {
        var elem = document.getElementById(elem_id);
        if (!elem) {
            return;
        }
        hc.JsonRenderElement(elem, elem_id);
    };

    hc.JsonRenderElement = function (elem, elem_id) {
        var entry = JSON.parse(elem.innerHTML);
        if (!entry) {
            return;
        }
        return hc.RenderElement(entry, elem_id);
    };

    hc._utilxArrayObjectHas = function (ar, v) {
        for (var i in ar) {
            if (ar[i] == v) {
                return true;
            }
        }
        return false;
    };

    hc.RenderUpdate = function (entry, elem_id) {
        if (entry.data.labels.length < 1) {
            return;
        }

        for (var i in entry.data.datasets) {
            if (entry.data.labels.length != entry.data.datasets[i].data.length) {
                return;
            }
        }

        var chart = hc.instances[elem_id];

        if (!chart) {
            return;
        }

        var plen = chart.data.labels.length;

        var offset = plen;
        for (var i in chart.data.labels) {
            if (chart.data.labels[i] == entry.data.labels[0]) {
                offset = i;
                break;
            }
        }

        if (offset < plen) {
            chart.data.labels.splice(offset);
            for (var i in chart.data.datasets) {
                chart.data.datasets[i].data.splice(offset);
            }
        }

        chart.data.labels = chart.data.labels.concat(entry.data.labels);

        for (var j in chart.data.datasets) {
            for (var k in entry.data.datasets) {
                if (
                    !entry.data.datasets[k].label ||
                    entry.data.datasets[k].label != chart.data.datasets[j].label
                ) {
                    continue;
                }

                chart.data.datasets[j].data = chart.data.datasets[j].data.concat(
                    entry.data.datasets[k].data
                );

                break;
            }
        }
        chart.update(200);

        setTimeout(function () {
            if (chart.data.labels.length > plen) {
                for (var j in chart.data.datasets) {
                    chart.data.datasets[j].data.splice(0, chart.data.labels.length - plen);
                }

                chart.data.labels.splice(0, chart.data.labels.length - plen);
            }
            chart.update(200);
        }, 300);
    };

    hc.RenderElement = function (entry, elem_id) {
        if (!entry || !entry.type) {
            return;
        }

        if (!entry.data || !entry.data.labels || !entry.data.datasets) {
            return;
        }

        entry.options = entry.options || {};
        entry.options.legend = entry.options.legend || {};

        entry.options.scales = entry.options.scales || {};
        entry.options.scales.x = entry.options.scales.x || {};
        entry.options.scales.y = entry.options.scales.y || {};

        entry.options.elements = entry.options.elements || {};

        entry.options.plugins = entry.options.plugins || {};
        entry.options.plugins.legend = entry.options.plugins.legend || {};


        if (entry.options.title && typeof entry.options.title === "string") {
            entry.options.plugins.title = {
                text: entry.options.title,
                display: true,
                fontSize: 16,
            };
        }

        if (!entry.options.legend.display) {
            entry.options.legend.display = true;
        } else {
            entry.options.legend.display = false;
        }

        // if (!entry.options.width) {
        //     entry.options.width = hc.opts_width;
        // }
        if (!entry.options.height) {
            entry.options.height = hc.opts_height;
        }

        entry.options.animation = {
            tension: {
                duration: 0,
                easing: 'linear',
            },
        };

        entry.options.hover = {
            animationDuration: 0,
        };

        entry.options.responsiveAnimationDuration = 0;

        entry.options.elements.line = {};

        for (var i in entry.data.datasets) {
            entry.data.datasets[i] = {
                // borderWidth: 0,
                label: entry.data.datasets[i].label,
                data: entry.data.datasets[i].data,
                tension: 0.1,
            };
        }

        entry.elem_id = elem_id;

        switch (entry.type) {
            case "bar-h":
            case "bar":
                hc.jsr_bar(entry);
                break;
            case "line":
                hc.jsr_line(entry);
                break;
            case "pie":
                hc.jsr_pie(entry);
                break;
        }
    };

    hc.jsr_line = function (entry) {
        entry.options.scales.x.grid = {
            display: false,
        }
        var cp = new color_pallet();
        for (var i in entry.data.datasets) {
            var clr = cp.rand();
            entry.data.datasets[i].borderColor = color.rgb(clr, 0.8);
            entry.data.datasets[i].backgroundColor = color.rgb(clr, 1);
            entry.data.datasets[i].fill = false;
        }
        if (entry.data.datasets.length == 1 && entry.options.title) {
            entry.options.legend.display = false;
        }
        if (!entry.options.plugins.legend.position) {
            entry.options.plugins.legend.position = "bottom";
        }
        hc.jsr_general(entry);
    };

    hc.jsr_bar = function (entry) {
        if (entry.type == "bar") {
            entry.options.scales.x.grid = {
                display: false,
            };
        }
        if (entry.type == "bar-h") {
            entry.options.scales.y.grid = {
                display: false,
            };
            entry.type = "horizontalBar";
        }

        if (entry.data.datasets.length == 1 && entry.options.title) {
            entry.options.legend.display = false;
        } else if (!entry.options.legend.position) {
            entry.options.legend.position = "bottom";
        }

        var cp = new color_pallet();
        for (var i in entry.data.datasets) {
            var clr = cp.rand();
            entry.data.datasets[i].stack = "stack-" + i;
            entry.data.datasets[i].borderColor = color.rgb(clr, 0.95);
            entry.data.datasets[i].backgroundColor = color.rgb(clr, 0.95);
        }

        if (entry.data.datasets.length == 1 && entry.options.title) {
            entry.options.legend.display = false;
        }

        hc.jsr_general(entry);
    };

    hc.jsr_pie = function (entry) {
        entry.options.scales.x = {};
        entry.options.scales.y = {};
        entry.options.cutoutPercentage = 50;

        if (!entry.options.legend.position) {
            entry.options.legend.position = "right";
        }

        for (var i in entry.data.datasets) {
            var cs = [];
            var cp = new color_pallet();
            for (var j = 0; j < entry.data.datasets[i].data.length; j++) {
                var clr = cp.rand();
                cs.push(color.rgb(clr, 1));
            }
            entry.data.datasets[i].backgroundColor = cs;
        }
        hc.jsr_general(entry);
    };

    hc.jsr_general = function (data) {
        hc.deps_load(function () {
            // Chart.defaults.global.title.fontSize = 16;
            var elem_c = document.createElement("canvas");
            elem_c.id = data.elem_id + "-hchart";
            if (data.options.width) {
                elem_c.width = "100vw";
            }
            if (data.options.height) {
                elem_c.height = "100vh";
            }

            var elem_d = document.createElement("div");
            if (data.options.width) {
                elem_d.style.width = data.options.width;
            }
            if (data.options.height) {
                elem_d.style.height = data.options.height;
            }
            elem_d.setAttribute("class", "hooto-chart-entry");
            elem_d.appendChild(elem_c);

            if (data.elem) {
                data.elem.insertAdjacentHTML("afterEnd", elem_d.outerHTML);
                data.elem.style.display = "none";
            } else {
                var elem = document.getElementById(data.elem_id);
                if (elem) {
                    elem.innerHTML = elem_d.outerHTML;
                }
            }

            var ctx = document.getElementById(data.elem_id + "-hchart").getContext("2d");
            var instance = new Chart(ctx, data);
            hc.instances[data.elem_id] = instance;
        });
    };
})(this);
