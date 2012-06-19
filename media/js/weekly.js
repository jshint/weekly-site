var editor;
var reportTemplate;

function onload() {
	var textarea = document.getElementById("editor");

	editor = CodeMirror.fromTextArea(textarea, {
		theme: 'default',
		lineNumbers: true
	});

	editor.focus();

	reportTemplate = document.getElementById("template-report").innerHTML;
}

function lint() {
	var linter = require("/src/jshint.js");
	var warnings = [];
	var messages = {}
	var table = document.querySelector(".editor .report table");
	var ret;

	try {
		ret = linter.lint({ code: editor.getValue() });
	} catch (err) {
		err = err.toString().split(":");

		warnings.push({
			code: "E000",
			line: err[1].replace("Line ", ""),
			desc: err[2]
		});
	}

	_.each(document.querySelectorAll(".editor .panel"), function (el) {
		el.style.display = "none";
	});

	if (ret) {
		if (_.size(ret.report.messages) === 0) {
			document.querySelector(".editor .success").style.display = "block";
			return;
		}

		messages = _.sortBy(ret.report.messages, function (msgs, line) { return line; });
		_.each(messages, function (msgs) {
			_.each(msgs, function (msg) {
				warnings.push({ code: msg.data.code, line: msg.line, desc: msg.data.desc });
			});
		});
	}

	table.innerHTML = _.template(reportTemplate, { warnings: warnings });
	document.querySelector(".editor .report").style.display = "block";
}
