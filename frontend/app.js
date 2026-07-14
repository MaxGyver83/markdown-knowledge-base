var API_BASE = window.API_BASE !== undefined ? window.API_BASE : 'http://localhost:8080';

let documents = [];
let currentId = null;

function $(sel) { return document.querySelector(sel); }
function $$(sel) { return document.querySelectorAll(sel); }

async function api(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.status === 204 ? null : res.json();
}

const API = {
  list: () => api('/documents'),
  get: (id) => api(`/documents/${id}`),
  create: (data) => api('/documents', { method: 'POST', body: JSON.stringify(data) }),
  update: (id, data) => api(`/documents/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (id) => api(`/documents/${id}`, { method: 'DELETE' }),
  reset: () => api('/documents/reset', { method: 'POST' }),
};

function navigate(hash) {
  window.location.hash = hash;
}

function toast(msg, type = 'success') {
  let el = $('#toast');
  if (!el) {
    el = document.createElement('div');
    el.id = 'toast';
    el.className = 'toast';
    document.body.appendChild(el);
  }
  el.textContent = msg;
  el.className = `toast ${type} show`;
  setTimeout(() => el.classList.remove('show'), 3000);
}

async function loadSidebar() {
  const list = $('#doc-list');
  const status = $('#sidebar-status');
  try {
    documents = await API.list();
    status.textContent = `${documents.length} document${documents.length !== 1 ? 's' : ''}`;
    list.innerHTML = documents.map(d => {
      const active = currentId === String(d.id) ? ' active' : '';
      return `<div class="doc-item${active}" data-id="${d.id}">${escapeHtml(d.title)}</div>`;
    }).join('');
    list.querySelectorAll('.doc-item').forEach(el => {
      el.addEventListener('click', () => navigate(`/view/${el.dataset.id}`));
    });
  } catch (e) {
    status.textContent = 'Connection error';
    list.innerHTML = `<div class="error">Could not load documents.<br>Is the backend running on ${API_BASE}?</div>`;
  }
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

/* -- Router -- */

function router() {
  const hash = window.location.hash.slice(1) || '/';
  const parts = hash.split('/').filter(Boolean);

  if (hash === '/' || parts[0] === '') {
    renderEmpty();
  } else if (parts[0] === 'view') {
    renderDetail(parts[1]);
  } else if (parts[0] === 'new') {
    renderForm(null);
  } else if (parts[0] === 'edit') {
    renderForm(parts[1]);
  } else {
    renderEmpty();
  }
}

/* -- Views -- */

function renderEmpty() {
  if (documents.length > 0) {
    navigate(`/view/${documents[0].id}`);
    return;
  }
  currentId = null;
  highlightSidebar();
  $('#view').innerHTML = `
    <div class="empty-state">
      <p>Select a document from the sidebar or create a new one.</p>
    </div>`;
}

async function renderDetail(id) {
  currentId = id;
  highlightSidebar();
  const view = $('#view');
  view.innerHTML = `<div class="spinner">Loading...</div>`;

  try {
    const doc = await API.get(id);
    const content = marked.parse(doc.content);
    view.innerHTML = `
      <div class="doc-header">
        <div class="doc-header-info">
          <h1>${escapeHtml(doc.title)}</h1>
          <div class="doc-meta">
            Created: ${formatDate(doc.created_at)}
            &middot;
            Updated: ${formatDate(doc.updated_at)}
          </div>
        </div>
        <div class="doc-actions">
          <button class="btn-secondary" onclick="navigate('/edit/${doc.id}')">Edit</button>
          <button class="btn-danger" id="btn-delete">Delete</button>
        </div>
      </div>
      <div class="doc-content">${content}</div>`;

    $('#btn-delete').addEventListener('click', () => handleDelete(doc.id, doc.title));
    $('#view').querySelectorAll('pre code:not(.language-mermaid)').forEach(b => hljs.highlightElement(b));
    document.querySelectorAll('pre code.language-mermaid').forEach(el => {
      const div = document.createElement('div');
      div.className = 'mermaid';
      div.textContent = el.textContent;
      el.parentElement.replaceWith(div);
    });
    if (document.querySelector('.mermaid')) {
      mermaid.run({ nodes: document.querySelectorAll('.mermaid') });
    }
  } catch (e) {
    view.innerHTML = `<div class="error">${escapeHtml(e.message)}</div>`;
  }
}

async function renderForm(id) {
  const view = $('#view');
  const isEdit = id !== null;

  view.innerHTML = `<div class="spinner">Loading...</div>`;
  currentId = isEdit ? null : null;
  highlightSidebar();

  let title = '';
  let content = '';

  if (isEdit) {
    try {
      const doc = await API.get(id);
      title = doc.title;
      content = doc.content;
      currentId = id;
      highlightSidebar();
    } catch (e) {
      view.innerHTML = `<div class="error">${escapeHtml(e.message)}</div>`;
      return;
    }
  }

  view.innerHTML = `
    <div class="form-container">
      <h2 style="margin-bottom:16px;font-size:24px;">${isEdit ? 'Edit Document' : 'New Document'}</h2>
      <div class="form-group">
        <label for="doc-title">Title</label>
        <input type="text" id="doc-title" value="${escapeHtml(title)}" placeholder="Document title" ${isEdit ? '' : 'autofocus'}>
      </div>
      <div class="form-group">
        <label for="doc-content">Content (Markdown)</label>
        <textarea id="doc-content" placeholder="Write your markdown here...">${escapeHtml(content)}</textarea>
      </div>
      <div class="form-actions">
        <button class="btn-primary" id="btn-save">${isEdit ? 'Update' : 'Create'}</button>
        <button class="btn-secondary" onclick="navigate('/')">Cancel</button>
      </div>
      <div id="form-error"></div>
    </div>`;

  $('#btn-save').addEventListener('click', () => handleSave(id));

  $('#doc-content').addEventListener('keydown', (e) => {
    if (e.ctrlKey && e.key === 'Enter') {
      e.preventDefault();
      handleSave(id);
    }
  });
}

function highlightSidebar() {
  $$('.doc-item').forEach(el => el.classList.toggle('active', el.dataset.id === currentId));
}

/* -- Actions -- */

async function handleSave(id) {
  const title = $('#doc-title').value.trim();
  const content = $('#doc-content').value;
  const errEl = $('#form-error');

  if (!title) {
    errEl.innerHTML = `<div class="error">Title is required.</div>`;
    return;
  }

  const btn = $('#btn-save');
  btn.disabled = true;
  btn.textContent = 'Saving...';

  try {
    if (id) {
      await API.update(id, { title, content });
      toast('Document updated');
    } else {
      const doc = await API.create({ title, content });
      toast('Document created');
      navigate(`/view/${doc.id}`);
    }
    await loadSidebar();
    if (id) {
      navigate(`/view/${id}`);
    }
  } catch (e) {
    errEl.innerHTML = `<div class="error">${escapeHtml(e.message)}</div>`;
    btn.disabled = false;
    btn.textContent = id ? 'Update' : 'Create';
  }
}

async function handleDelete(id, title) {
  if (!confirm(`Delete "${title}"?`)) return;

  try {
    await API.delete(id);
    toast('Document deleted');
    await loadSidebar();
    navigate('/');
  } catch (e) {
    toast(e.message, 'error');
  }
}

async function handleReset() {
  if (!confirm('Reset all documents? This will delete all documents and restore the demo documents.')) return;

  try {
    await API.reset();
    toast('Demo documents restored');
    await loadSidebar();
    navigate('/');
  } catch (e) {
    toast(e.message, 'error');
  }
}

function formatDate(iso) {
  const d = new Date(iso);
  return d.toLocaleDateString('en-US', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  });
}

/* -- Init -- */

document.addEventListener('DOMContentLoaded', () => {
  $('#btn-new').addEventListener('click', () => navigate('/new'));
  $('#btn-reset').addEventListener('click', handleReset);
  window.addEventListener('hashchange', router);
  loadSidebar().then(router);
});
