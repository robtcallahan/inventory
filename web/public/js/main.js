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
    delete: document.getElementById('item-delete'),
    cancel: document.getElementById('item-cancel'),
    save: document.getElementById('item-save'),
    confirmDialog: document.getElementById('delete-confirm-dialog'),
    confirmMessage: document.getElementById('delete-confirm-message'),
    confirmDelete: document.getElementById('confirm-delete'),
    cancelDelete: document.getElementById('cancel-delete'),
    photoGallery: document.getElementById('item-photo-gallery'),
    photoInput: document.getElementById('item-photo-input'),
    photoUploadLabel: document.getElementById('item-photo-upload-label'),
    photoStatus: document.getElementById('item-photo-status'),
    photoViewer: document.getElementById('photo-viewer'),
    photoViewerImage: document.getElementById('photo-viewer-image'),
    photoViewerClose: document.getElementById('photo-viewer-close'),
};

const tableFields = ['number', 'quantity', 'name', 'short_description', 'replacement_value', 'resale_value',
    'total_value', 'manufacturer', 'model', 'category', 'purchase_date', 'serial_number'];
const detailFields = ['quantity', 'name', 'short_description', 'replacement_value', 'resale_value', 'total_value',
    'manufacturer', 'model', 'category', 'purchase_date', 'serial_number', 'details'];
const numericFields = new Set(['quantity', 'replacement_value', 'resale_value']);
const moneyFields = new Set(['replacement_value', 'resale_value', 'total_value']);

let items = [];
let selectedItem = null;
let selectedCategory = '';
let searchText = "";
let saving = false;
let deleting = false;

function fieldLabel(key) {
    if (key === 'quantity') return 'Qty';
    if (key === 'number') return 'Num';
    if (key === 'short_description') return 'Short Description';
    if (key === 'replacement_value') return 'Replacement Value';
    if (key === 'resale_value') return 'Resale Value';
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
            item.short_description,
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
        a.short_description.localeCompare(b.short_description)
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
    ui.replacement_value.textContent = money(displayed.reduce((sum, item) => sum + Number(item.total_value || 0), 0));
}

function createFields() {
    const category = ui.filter.cloneNode(true);
    category.remove(0);  // Remove "All"
    category.id = 'item-field-category';
    category.name = 'category';

    for (const key of detailFields) {
        const field = document.createElement('div');
        field.className = 'item-fields__row';
        const label = document.createElement('label');
        label.className = 'item-fields__label';
        label.htmlFor = `item-field-${key}`;
        label.textContent = fieldLabel(key);
        let input = key;
        if (key === "category") {
            input = category;
        } else if (key === "details") {
            input = document.createElement("textarea");
            input.rows = 3;
        } else {
            input = document.createElement("input");
            input.type = numericFields.has(key) ? 'number' : 'text';
        }
        input.className = 'item-fields__input';
        input.id = `item-field-${key}`;
        input.name = key;
        // if (key !== 'category') {
        //     input.type = numericFields.has(key) ? 'number' : 'text';
        // }
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
    for (const input of ui.fields.querySelectorAll('input, select, textarea')) {
        input.disabled = !editing || input.name === 'total_value';
    }
    ui.close.classList.toggle('is-hidden', editing);
    ui.edit.classList.toggle('is-hidden', editing);
    ui.delete.classList.toggle('is-hidden', editing);
    ui.cancel.classList.toggle('is-hidden', !editing);
    ui.save.classList.toggle('is-hidden', !editing);
    ui.photoUploadLabel.classList.toggle('is-hidden', !editing || adding || !selectedItem?.id);
    ui.photoGallery.querySelectorAll('.photo-delete').forEach(button => button.classList.toggle('is-hidden', !editing));
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
        ui.form.elements.namedItem(key).replacement_value = item[key] ?? '';
    }
    setMode(false);
    loadPhotos(item.id);
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
    ui.photoGallery.replaceChildren();
    ui.photoStatus.textContent = 'Save the item before adding photos.';
    setMode(true, true);
    ui.dialog.classList.remove('is-hidden');
    ui.edit.focus();
    ui.form.elements.namedItem('name').focus();
}

function readItem() {
    const item = { id: selectedItem?.id ?? 0 };
    for (const key of detailFields) {
        if (key === 'total_value') continue; // Computed by the server.
        const value = ui.form.elements.namedItem(key).replacement_value.trim();
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
    ui.delete.disabled = true;
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
        closeItemAfterSaveOrDelete();
    } catch (error) {
        showError(error.message);
    } finally {
        saving = false;
        ui.save.disabled = false;
        ui.delete.disabled = false;
    }
}

async function deleteItem(event) {
    event.preventDefault();
    if (deleting || !selectedItem) return;
    const payload = { id: selectedItem.id };
    deleting = true;
    ui.save.disabled = true;
    ui.delete.disabled = true;
    ui.confirmDelete.disabled = true;
    showError();
    try {
        const response = await fetch('/api/delete', {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!response.ok) throw new Error(`Delete failed (HTTP ${response.status}).`);
        // Fetch the server's authoritative computed fields and formatted dates.
        await loadItems();
        ui.confirmDialog.close();
        closeItemAfterSaveOrDelete();
    } catch (error) {
        showError(error.message);
    } finally {
        deleting = false;
        ui.save.disabled = false;
        ui.delete.disabled = false;
        ui.confirmDelete.disabled = false;
    }
}

function confirmDeleteItem() {
    if (saving || deleting || !selectedItem) return;

    ui.confirmMessage.textContent =
        `Are you sure you want to delete "${selectedItem.name}"?`;

    ui.confirmDialog.showModal();
    ui.cancelDelete.focus();
}

function closeItemAfterSaveOrDelete() {
    ui.dialog.classList.add('is-hidden');
    selectedItem = null;
}

function searchFilter(event) {
    searchText = event.target.replacement_value.trim().toLowerCase();
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
        selectedCategory = ui.filter.replacement_value;
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
    ui.photoInput.addEventListener('change', uploadSelectedPhoto);
    setupPhotoDropZone();
    ui.photoViewerClose.addEventListener('click', () => ui.photoViewer.close());
    ui.delete.addEventListener('click', confirmDeleteItem);
    ui.cancelDelete.addEventListener('click', () => {
        ui.confirmDialog.close();
    });
    ui.confirmDelete.addEventListener('click', deleteItem);
    ui.dialog.addEventListener('keydown', event => {
        if (event.key === 'Escape' && !saving) closeItem();
    });
    loadItems().catch(error => {
        console.error(error);
        ui.body.textContent = 'Unable to load inventory. Please refresh the page.';
    });
}

init();

async function loadPhotos(itemID) {
    ui.photoGallery.replaceChildren();
    ui.photoStatus.textContent = '';
    if (!itemID) return;
    try {
        const response = await fetch(`/api/items/${itemID}/photos`);
        if (!response.ok) throw new Error('Unable to load photos.');
        const photos = await response.json();
        if (selectedItem?.id !== itemID) return;
        if (!photos.length) ui.photoStatus.textContent = 'No photos yet.';
        for (const photo of photos) {
            const tile = document.createElement('div');
            tile.className = 'photo-tile';
            const url = `/api/items/${itemID}/photos/${encodeURIComponent(photo.filename)}`;
            const view = document.createElement('button');
            view.type = 'button';
            view.className = 'photo-view';
            const img = document.createElement('img');
            img.src = url;
            img.alt = 'Inventory photo';
            img.loading = 'lazy';
            view.append(img);
            view.addEventListener('click', () => {
                ui.photoViewerImage.src = url;
                ui.photoViewer.showModal();
            });
            const remove = document.createElement('button');
            remove.type = 'button';
            remove.className = 'button button--danger photo-delete';
            remove.textContent = 'Delete';
            remove.classList.toggle('is-hidden', ui.save.classList.contains('is-hidden'));
            remove.addEventListener('click', async () => {
                if (!confirm('Delete this photo?')) return;
                remove.disabled = true;
                try {
                    const result = await fetch(`/api/items/${itemID}/photos/${photo.id}`, {method:'DELETE'});
                    if (!result.ok) throw new Error('Unable to delete photo.');
                    await loadPhotos(itemID);
                } catch (error) {
                    ui.photoStatus.textContent = error.message;
                    remove.disabled = false;
                }
            });
            tile.append(view, remove);
            ui.photoGallery.append(tile);
        }
    } catch (error) {
        ui.photoStatus.textContent = error.message;
    }
}

async function uploadSelectedPhoto() {
    const file = ui.photoInput.files[0];
    await uploadPhotoFile(file);
}

async function uploadPhotoFile(file) {
    const itemID = selectedItem?.id;
    if (!file || !itemID) return;
    if (file.size > 10 * 1024 * 1024) {
        ui.photoStatus.textContent = 'Photos must be 10 MB or smaller.';
        ui.photoInput.replacement_value = '';
        return;
    }
    const form = new FormData();
    form.append('photo', file);
    ui.photoInput.disabled = true;
    ui.photoStatus.textContent = 'Uploading photo…';
    try {
        const response = await fetch(`/api/items/${itemID}/photos`, {method:'POST',body:form});
        if (!response.ok) throw new Error(await response.text() || 'Upload failed.');
        await loadPhotos(itemID);
    } catch (error) {
        ui.photoStatus.textContent = error.message;
    } finally {
        ui.photoInput.replacement_value = '';
        ui.photoInput.disabled = false;
    }
}

function setupPhotoDropZone() {
    const zone = document.getElementById('item-photos');
    let dragDepth = 0;
    const canUpload = () => selectedItem?.id > 0 &&
        !ui.photoUploadLabel.classList.contains('is-hidden') &&
        !ui.photoInput.disabled;
    const isFileDrag = event => Array.from(event.dataTransfer?.types || []).includes('Files');

    zone.addEventListener('dragenter', event => {
        if (!isFileDrag(event) || !canUpload()) return;
        event.preventDefault();
        dragDepth++;
        zone.classList.add('is-dragging');
    });
    zone.addEventListener('dragover', event => {
        if (!isFileDrag(event) || !canUpload()) return;
        event.preventDefault();
        event.dataTransfer.dropEffect = 'copy';
    });
    zone.addEventListener('dragleave', event => {
        if (!isFileDrag(event)) return;
        dragDepth = Math.max(0, dragDepth - 1);
        if (!dragDepth) zone.classList.remove('is-dragging');
    });
    zone.addEventListener('drop', async event => {
        if (!isFileDrag(event)) return;
        event.preventDefault(); // Never navigate away to the dropped image.
        dragDepth = 0;
        zone.classList.remove('is-dragging');
        if (!canUpload()) {
            ui.photoStatus.textContent = 'Click Edit to add a photo.';
            return;
        }
        const files = Array.from(event.dataTransfer.files);
        if (!files.length) {
            ui.photoStatus.textContent = 'No file received. Try dragging from Finder instead.';
            return;
        }
        if (files.length !== 1) {
            ui.photoStatus.textContent = 'Please drop one photo at a time.';
            return;
        }
        const file = files[0];
        if (!['image/jpeg', 'image/png', 'image/gif', 'image/webp'].includes(file.type)) {
            ui.photoStatus.textContent = 'Please drop a JPEG, PNG, GIF, or WebP photo.';
            return;
        }
        await uploadPhotoFile(file);
    });
    // Keep an image dropped outside the gallery from replacing the inventory page.
    document.addEventListener('dragover', event => {
        if (isFileDrag(event) && !ui.dialog.classList.contains('is-hidden')) event.preventDefault();
    });
    document.addEventListener('drop', event => {
        if (isFileDrag(event) && !ui.dialog.classList.contains('is-hidden')) event.preventDefault();
    });
}
