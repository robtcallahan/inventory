'use strict';

// DOM hooks use IDs. CSS classes describe appearance; is-* classes describe state.
const ui = {
    filter: document.getElementById('category-filter'),
    search: document.getElementById('search-filter'),
    count: document.getElementById('inventory-count'),
    value: document.getElementById('inventory-value'),
    add: document.getElementById('inventory-add'),
    head: document.getElementById('inventory-table-head'),
    body: document.getElementById('inventory-table-body'),
    dialog: document.getElementById('item-dialog'),
    dialogTitle: document.getElementById('item-dialog-title'),
    form: document.getElementById('item-form'),
    fields: document.getElementById('item-fields'),
    error: document.getElementById('item-form-error'),
    close: document.getElementById('item-close'),
    edit: document.getElementById('item-edit'),
    cancel: document.getElementById('item-cancel'),
    save: document.getElementById('item-save'),
};

const tableFields = ['number', 'quantity', 'name', 'description', 'value',
    'total_value', 'manufacturer', 'model', 'category', 'purchase_date', 'serial_number'];
const detailFields = ['quantity', 'name', 'description', 'value', 'total_value',
    'manufacturer', 'model', 'category', 'purchase_date', 'serial_number', 'details'];
const numericFields = new Set(['quantity', 'value']);
const moneyFields = new Set(['value', 'total_value']);

let items = [];
let selectedItem = null;
let selectedCategory = '';
let searchText = "";
let saving = false;

function fieldLabel(key) {
    if (key === 'quantity') return 'Qty';
    if (key === 'number') return 'Num';
    return key.replaceAll('_', ' ').replace(/\b\w/g, char => char.toUpperCase());
}

function money(amount) {
    return '$' + Number(amount || 0).toLocaleString('en-US');
}


function visibleItems() {
    return items.filter(item => {
        const matchesCategory =
            !selectedCategory || item.category === selectedCategory;

        const matchesSearch = [
            item.name,
            item.description,
            item.manufacturer,
            item.model,
            item.serial_number,
            item.details
        ].some(value =>
            String(value ?? "").toLowerCase().includes(searchText)
        );

        return matchesCategory && matchesSearch;
    }).sort((a, b) =>
        a.name.localeCompare(b.name) ||
        a.description.localeCompare(b.description)
    );
}

function renderTable() {
    const displayed = visibleItems();
    const header = document.createElement('tr');
    for (const key of tableFields) {
        const th = document.createElement('th');
        th.scope = 'col';
        th.textContent = fieldLabel(key);
        header.append(th);
    }
    ui.head.replaceChildren(header);
    ui.body.replaceChildren();

    displayed.forEach((item, index) => {
        const row = document.createElement('tr');
        row.tabIndex = 0;
        row.setAttribute('aria-label', `View ${item.name}`);
        for (const key of tableFields) {
            const td = document.createElement('td');
            if (moneyFields.has(key)) {
                td.classList.add('inventory-table__money');
                td.textContent = money(item[key]);
            } else {
                td.textContent = String(key === 'number' ? index + 1 : (item[key] ?? ''));
            }
            row.append(td);
        }
        row.addEventListener('click', () => openItem(item));
        row.addEventListener('keydown', event => {
            if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault();
                openItem(item);
            }
        });
        ui.body.append(row);
    });

    ui.count.textContent = displayed.length;
    ui.value.textContent = money(displayed.reduce((sum, item) => sum + Number(item.total_value || 0), 0));
}

function createFields() {
    for (const key of detailFields) {
        const field = document.createElement('div');
        field.className = 'item-fields__row';
        const label = document.createElement('label');
        label.className = 'item-fields__label';
        label.htmlFor = `item-field-${key}`;
        label.textContent = fieldLabel(key);
        const input = document.createElement('input');
        input.className = 'item-fields__input';
        input.id = `item-field-${key}`;
        input.name = key;
        input.type = numericFields.has(key) ? 'number' : 'text';
        if (numericFields.has(key)) {
            input.step = '1';
            input.min = '0';
        }
        if (key === 'total_value') input.readOnly = true;
        field.append(label, input);
        ui.fields.append(field);
    }
}

function showError(message = '') {
    ui.error.textContent = message;
    ui.error.classList.toggle('is-hidden', !message);
}

function setMode(editing, adding = false) {
    for (const input of ui.fields.querySelectorAll('input')) {
        input.disabled = !editing || input.name === 'total_value';
    }
    ui.close.classList.toggle('is-hidden', editing);
    ui.edit.classList.toggle('is-hidden', editing);
    ui.cancel.classList.toggle('is-hidden', !editing);
    ui.save.classList.toggle('is-hidden', !editing);
    if (adding) {
        ui.dialogTitle.textContent = "Add Item";
    } else if (editing) {
        ui.dialogTitle.textContent = "Edit Item";
    } else {
        ui.dialogTitle.textContent = "Item Details";
    }}

function openItem(item) {
    selectedItem = item;
    showError();
    for (const key of detailFields) {
        ui.form.elements.namedItem(key).value = item[key] ?? '';
    }
    setMode(false);
    ui.dialog.classList.remove('is-hidden');
    ui.edit.focus();
}

function closeItem() {
    if (saving) return;
    ui.dialog.classList.add('is-hidden');
    selectedItem = null;
}

function addItem(event) {
    selectedItem = {id: 0};
    setMode(true, true);
    ui.dialog.classList.remove('is-hidden');
    ui.edit.focus();
    ui.form.elements.namedItem('name').focus();
}

function readItem() {
    const item = { id: selectedItem?.id ?? 0 };
    for (const key of detailFields) {
        if (key === 'total_value') continue; // Computed by the server.
        const value = ui.form.elements.namedItem(key).value.trim();
        if (numericFields.has(key)) {
            if (!/^\d+$/.test(value)) throw new Error(`${fieldLabel(key)} must be a nonnegative whole number.`);
            item[key] = Number(value);
            if (!Number.isSafeInteger(item[key])) throw new Error(`${fieldLabel(key)} is too large.`);
        } else if (key === 'purchase_date') {
            if (!value) {
                item[key] = null;
            } else {
                const match = /^(\d{2})\/(\d{2})\/(\d{4})$/.exec(value);
                if (!match) throw new Error('Purchase Date must use MM/DD/YYYY.');
                const [, mm, dd, yyyy] = match.map(Number);
                const date = new Date(Date.UTC(yyyy, mm - 1, dd));
                if (date.getUTCFullYear() !== yyyy || date.getUTCMonth() !== mm - 1 || date.getUTCDate() !== dd) {
                    throw new Error('Purchase Date is invalid.');
                }
                item[key] = date.toISOString();
            }
        } else {
            item[key] = value;
        }
    }
    return item;
}

async function saveItem(event) {
    event.preventDefault();
    if (saving || !selectedItem) return;
    let payload;
    try {
        payload = readItem();
    } catch (error) {
        showError(error.message);
        return;
    }
    saving = true;
    ui.save.disabled = true;
    showError();
    try {
        const response = await fetch('/api/save', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!response.ok) throw new Error(`Save failed (HTTP ${response.status}).`);
        // Fetch the server's authoritative computed fields and formatted dates.
        await loadItems();
        closeItemAfterSave();
    } catch (error) {
        showError(error.message);
    } finally {
        saving = false;
        ui.save.disabled = false;
    }
}

function closeItemAfterSave() {
    ui.dialog.classList.add('is-hidden');
    selectedItem = null;
}

function searchFilter(event) {
    searchText = event.target.value.trim().toLowerCase();
    renderTable();
}

async function loadItems() {
    const response = await fetch('/api/items');
    if (!response.ok) throw new Error(`Unable to load items (HTTP ${response.status}).`);
    items = await response.json();
    renderTable();
}

function init() {
    createFields();
    ui.filter.addEventListener('change', () => {
        selectedCategory = ui.filter.value;
        renderTable();
    });
    ui.search.addEventListener('input', searchFilter)
    ui.add.addEventListener('click', addItem);
    ui.close.addEventListener('click', closeItem);
    ui.cancel.addEventListener('click', closeItem);
    ui.edit.addEventListener('click', () => {
        setMode(true);
        ui.form.elements.namedItem('name').focus();
    });
    ui.form.addEventListener('submit', saveItem);
    ui.dialog.addEventListener('keydown', event => {
        if (event.key === 'Escape' && !saving) closeItem();
    });
    loadItems().catch(error => {
        console.error(error);
        ui.body.textContent = 'Unable to load inventory. Please refresh the page.';
    });
}

init();
