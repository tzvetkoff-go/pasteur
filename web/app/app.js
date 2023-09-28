// Monaco
import * as monaco from 'monaco-editor/esm/vs/editor/editor.main.js';

// jQuery & select2
import jquery from 'jquery'
import select2 from 'select2';
import '../../node_modules/select2/dist/css/select2.css'
import '../select2-dark-theme/select2-dark-theme.scss'

// Stylesheets
import './app.scss';

// Plug select2 into jQuery
select2(jquery.prototype);

//
// Application
//

document.addEventListener('DOMContentLoaded', () => {
  // Relative URL root
  if (typeof window.RelativeURLRoot === 'undefined') {
    window.RelativeURLRoot = '';
  }

  // Editor theme
  monaco.editor.defineTheme('pasteur', {
    base: 'vs-dark',
    inherit: true,
    rules: [{
      background: '#04121b',
    }],
    colors: {
      'editor.background': '#04121b',
      'editor.foreground': '#d4dbdf',
      'editor.inactiveSelectionBackground': '#14507d',
      'editor.lineHighlightBackground': '#114368',
      'editor.selectionBackground': '#175d92',
      'editor.selectionForeground': '#175d92',
      'editorLineNumber.background': '#04121b',
      'editorLineNumber.foreground': '#677e8d',
      'editorWidget.background': '#0c1a24',
      'editorWidget.border': '#082437',
      'editorRuler.foreground': '#082437',
      'input.background': '#082437',
      'input.border': '#1d4057',
      'input.foreground': '#ffffff',
      'scrollbar.shadow': '#00000060',
      'progressBar.background': '#4793cc'
    }
  });

  // Editor environment
  window.MonacoEnvironment = {
    getWorkerUrl: (moduleId, label) => {
      if (label === 'json') {
        return window.RelativeURLRoot + '/assets/vs/language/json/json.worker.js';
      }
      if (label === 'css' || label === 'scss' || label === 'less') {
        return window.RelativeURLRoot + '/assets/vs/language/css/css.worker.js';
      }
      if (label === 'html' || label === 'handlebars' || label === 'razor') {
        return window.RelativeURLRoot + '/assets/vs/language/html/html.worker.js';
      }
      if (label === 'typescript' || label === 'javascript') {
        return window.RelativeURLRoot + '/assets/vs/language/typescript/ts.worker.js';
      }

      return window.RelativeURLRoot + '/assets/vs/editor/editor.worker.js';
    }
  };

  // Editor initialization
  const codeEditor = document.getElementById('content');
  let monacoEditor = null;
  if (codeEditor !== null) {
    monacoEditor = monaco.editor.create(codeEditor.parentNode, {
      theme: 'pasteur',
      fontSize: 14,
      rulers: [80, 100, 120],
      overviewRulerLanes: 0,
      minimap: {
        enabled: false
      },
      renderLineHighlight: 'all',
      renderLineHighlightOnlyWhenFocus: false,
      renderWhitespace: 'all',
      scrollbar: {
        horizontalScrollbarSize: 6,
        verticalScrollbarSize: 6
      },
      scrollBeyondLastLine: false,
      value: codeEditor.value,
      insertSpaces: true,
      tabSize: 4,
      readOnly: codeEditor.readOnly,
      language: codeEditor.dataset.language,
      wordWrap: 'on'
    });

    monacoEditor.addAction({
      id: 'toggle-wrap',
      label: 'Toggle Word Wrap',
      run: (editor) => {
        editor.updateOptions({
          wordWrap: editor.getRawOptions().wordWrap == 'on' ? 'off' : 'on'
        });
      }
    });

    monacoEditor.getModel().onDidChangeContent(() => {
      codeEditor.value = monacoEditor.getValue();
      codeEditor.dispatchEvent(new Event('change'));
    });

    window.addEventListener('resize', () => {
      monacoEditor.layout();
    });
  }

  // Select initialization
  [].forEach.call(document.getElementsByTagName('select'), (select) => {
    if (select.nextSibling !== undefined && select.nextSibling.nodeType === 3) {
      select.parentNode.removeChild(select.nextSibling);
    }

    const options = {
      theme: 'dark'
    };

    const firstOption = select.getElementsByTagName('option')[0];
    if (firstOption !== undefined && firstOption.value === '') {
      options.placeholder = firstOption.innerText;
      options.allowClear = true;
    };

    jquery(select).select2(options);
  });

  // jQuery 3.6.0 fix
  jquery(document).on('select2:open', () => {
    document.querySelector('.select2-search__field', this).focus();
  });

  // Language
  const filetypeSelect = document.getElementById('filetype')
  if (filetypeSelect !== null && monacoEditor !== null) {
    filetypeSelect.onchange = () => {
      monaco.editor.setModelLanguage(monacoEditor.getModel(), filetypeSelect.value);
    };
  }

  // Indent style
  const indentStyleSelect = document.getElementById('indent-style');
  if (indentStyleSelect !== null && monacoEditor !== null) {
    indentStyleSelect.onchange = () => {
      monacoEditor.getModel().updateOptions({
        insertSpaces: indentStyleSelect.value === 'spaces'
      });
    };
  }

  // Indent size
  const indentSizeSelect = document.getElementById('indent-size');
  if (indentSizeSelect !== null && monacoEditor !== null) {
    indentSizeSelect.onchange = () => {
      monacoEditor.getModel().updateOptions({
        tabSize: parseInt(indentSizeSelect.value, 10)
      });
    };
  }

  // Submit paste
  const submitPasteButton = document.getElementById('submit-paste');
  if (submitPasteButton !== null) {
    submitPasteButton.onclick = () => {
      const content = document.getElementById('content');
      if (content.value !== '') {
        return true;
      }

      alert('Content cannot be empty!');
      return false;
    };
  }

  // Clone paste
  const clonePasteButton = document.getElementById('clone-paste');
  if (clonePasteButton !== null) {
    clonePasteButton.onclick = () => {
      clonePasteButton.style.display = 'none';
      submitPasteButton.style.display = 'inline-block';
      jquery('select').prop('disabled', false);
      monacoEditor.updateOptions({
        readOnly: false
      });
    };
  }
});
