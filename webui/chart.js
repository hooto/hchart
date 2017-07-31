(function(global, undefined) {
'use strict';

if (global.hooto_chart) {
    return
}

var chart = global.hooto_chart = {
    version: "0.1.0.dev",
    basepath: "/chart/~",
    opts_width: "800px",
    opts_height: "600px",
}

var color = {
    theme: null,
    theme_tb: ["337ab7", "5cb85c", "5bc0de", "f0ad4e", "d9534f", "333333"],
    theme_gc: ["0057e7", "d62d20", "ffa700", "008744", "9c27b0", "333333"],
}

function color_pallet(theme)
{
    if (!theme) {
        this.theme = color.theme_gc;    
    }
    this.randlast  = 0;
    this.reset  = 1;

    this.rand = function() {
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
			color_i16 = ((1 << 24) + (r << 16) + (g << 8) + b);
		}
		return color_i16;
    }
}

color.rgb = function(color_i16, alpha)
{
    var r = (color_i16 >> 16) & 255;
    var g = (color_i16 >> 8) & 255;
    var b = color_i16 & 255;

    return "rgba("+ r +","+ g +","+ b +","+ alpha +")";
}

chart.deps_load = function(cb)
{
    seajs.use([
        chart.basepath + "/chartjs/chart.js",
    ], cb);
}

chart.JsonRenderElementID = function(elem_id)
{
    var elem = document.getElementById(elem_id);
    if (!elem) {
        return;
    }
    chart.JsonRenderElement(elem, elem_id);
}

chart.JsonRenderElement = function(elem, elem_id)
{
    var entry = JSON.parse(elem.innerHTML); 
    if (!entry || !entry.type) {
        return;
    }

    if (!entry.data || !entry.data.labels || !entry.data.datasets) {
        return;
    }

    entry.options = entry.options || {};
    entry.options.scales = entry.options.scales || {};
    entry.options.elements = entry.options.elements || {};

	if (!entry.options.width) {
		entry.options.width = chart.opts_width;
	}
	if (!entry.options.height) {
		entry.options.height = chart.opts_height;
	}

    entry.options.animation = {
        duration: 0,
    };

    entry.options.hover = {
        animationDuration: 0,
    };

    entry.options.responsiveAnimationDuration = 0;

    entry.options.elements.line = {
    };

    entry.options.scales.xAxes = [{
        // position: 'top', 
    }];
    entry.options.scales.yAxes = [{
        // stacked: true,
    }];

    for (var i in entry.data.datasets) {
        entry.data.datasets[i] = {
            borderWidth: 0,
            label: entry.data.datasets[i].label,
            data: entry.data.datasets[i].data,
        };
    }

    entry.elem = elem;
    entry.elem_id = elem_id;

    switch (entry.type) {
    case "bar-h":
    case "bar":
        chart.jsr_bar(entry);
        break;
    case "line":
        chart.jsr_line(entry);
        break;
    case "pie":
        chart.jsr_pie(entry);
        break;
    }
}

chart.jsr_line = function(entry)
{
    entry.options.scales.xAxes[0].gridLines = {
        display: false,
    };
    var cp = new color_pallet();
    for (var i in entry.data.datasets) {
        var clr = cp.rand();
        entry.data.datasets[i].borderColor = color.rgb(clr, 0.8);
        entry.data.datasets[i].backgroundColor = color.rgb(clr, 1);
        entry.data.datasets[i].fill = false;
    }
    chart.jsr_general(entry);
}

chart.jsr_bar = function(entry)
{
    if (entry.type == "bar") {
        entry.options.scales.xAxes[0].gridLines = {
            display: false,
        };
    }
    if (entry.type == "bar-h") {
        entry.options.scales.yAxes[0].gridLines = {
            display: false,
        };
        entry.type = "horizontalBar";
    }

    var cp = new color_pallet();
    for (var i in entry.data.datasets) {
        var clr = cp.rand();
        entry.data.datasets[i].stack = "stack-" + i;
        entry.data.datasets[i].borderColor = color.rgb(clr, 0.95);
        entry.data.datasets[i].backgroundColor = color.rgb(clr, 0.95);
    }

    chart.jsr_general(entry);
}

chart.jsr_pie = function(entry)
{
    entry.options.scales.xAxes = {};
    entry.options.scales.yAxes = {};
    entry.options.cutoutPercentage = 50;

    for (var i in entry.data.datasets) {
        var cs = [];
        var cp = new color_pallet();
        for (var j = 0; j < entry.data.datasets[i].data.length; j++) {
            var clr = cp.rand();
            cs.push(color.rgb(clr, 1));
        }
        entry.data.datasets[i].backgroundColor = cs;
    }
    chart.jsr_general(entry);
}

chart.jsr_general = function(data)
{
    chart.deps_load(function(){

        var elem_c = document.createElement("canvas");
        elem_c.id = data.elem_id + "-chart";
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

        data.elem.insertAdjacentHTML("afterEnd", elem_d.outerHTML);
        data.elem.style.display = "none";

        var ctx = document.getElementById(data.elem_id +"-chart").getContext("2d");
        new Chart(ctx, data);
    });
}

})(this);
