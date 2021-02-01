//
// Application
//

document.addEventListener('DOMContentLoaded', function() {
  // Control wrappers
  var newPasteMeta = document.getElementsByClassName('new-paste-meta')[0];
  var editPasteMeta = document.getElementsByClassName('edit-paste-meta')[0];

  // The paste form
  var pasteForm = document.getElementsByClassName('paste-form')[0];

  // Indent style select
  var indentStyleSelect = document.getElementsByClassName('indent-style-select')[0];
  jQuery(indentStyleSelect).select2({
    theme: 'dark'
  })

  // Indent size select
  var indentSizeSelect = document.getElementsByClassName('indent-size-select')[0];
  jQuery(indentSizeSelect).select2({
    theme: 'dark'
  });

  // Mime type select
  var mimeTypeSelect = document.getElementsByClassName('mime-type-select')[0];
  jQuery(mimeTypeSelect).select2({
    theme: 'dark'
  });

  // Code editor
  var codeEditor = document.getElementsByClassName('code-editor')[0];
  var codeMirror = CodeMirror.fromTextArea(codeEditor, {
    lineNumbers: true,
    styleActiveLine: {
      nonEmpty: true
    },
    matchBrackets: true,
    indentWithTabs: indentStyleSelect.value == 'tabs',
    indentUnit: parseInt(indentSizeSelect.value, 10),
    tabSize: parseInt(indentSizeSelect.value, 10),
    mode: mimeTypeSelect.value,
    readOnly: codeEditor.readOnly,
    theme: 'material-darker',
    showWhitespaces: true,
    extraKeys: {
      'Tab': function(cm) {
        if (cm.somethingSelected()) {
          cm.execCommand('indentMore');
        } else {
          cm.execCommand(indentStyleSelect.value == 'tabs' ? 'insertTab' : 'insertSoftTab');
        }
      },
      'Shift-Tab': function(cm) {
        cm.execCommand('indentLess');
      },
      'Ctrl-Enter': function(cm) {
        pasteForm.submit();
      }
    }
  });

  // Indent style change
  indentStyleSelect.onchange = function() {
    codeMirror.setOption('indentWithTabs', indentStyleSelect.value == 'tabs');
  };

  // Indent size change
  indentSizeSelect.onchange = function() {
    codeMirror.setOption('indentUnit', parseInt(indentSizeSelect.value, 10));
    codeMirror.setOption('tabSize', parseInt(indentSizeSelect.value, 10));
  };

  // Mime type change
  mimeTypeSelect.onchange = function() {
    codeMirror.setOption('mode', mimeTypeSelect.value);
  };

  // Edit paste
  var editPasteButton = document.getElementsByClassName('edit-paste-button')[0];
  editPasteButton.onclick = function() {
    editPasteMeta.style.visibility = 'hidden';
    newPasteMeta.style.visibility = 'visible';
    codeMirror.setOption('readOnly', false);
    codeMirror.focus();
  };

  // More edit
  if (codeEditor.readOnly) {
    newPasteMeta.style.visibility = 'hidden';
    editPasteMeta.style.visibility = 'visible';
  } else {
    codeMirror.focus();
  }
});
