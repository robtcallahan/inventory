'use strict';

// DOM hooks use IDs. CSS classes describe appearance; is-* classes describe state.
const ui = {
    filterButton: document.getElementById('inventory-filter'),
    toolbarFilters: document.getElementById('toolbar-filters'),
    filtersCount: document.getElementById('filters-count'),

    categoryFilter: document.getElementById('category-filter'),
    photoFilter: document.getElementById('photo-filter'),
    infoFilter: document.getElementById('info-filter'),
    roomFilter: document.getElementById('room-filter'),
    clearFilters: document.getElementById('clear-filters'),

    bulk: document.getElementById('inventory-bulk'),
    bulkDialog: document.getElementById('bulk-dialog'),
    bulkCount: document.getElementById('bulk-count'),
    bulkField: document.getElementById('bulk-field'),
    bulkValue: document.getElementById('bulk-value'),
    bulkError: document.getElementById('bulk-error'),
    bulkApply: document.getElementById('bulk-apply'),
    bulkCancel: document.getElementById('bulk-cancel'),

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
    'total_value', 'manufacturer', 'model', 'category', 'room', 'info_complete', 'purchase_date', 'serial_number'];
const detailFields = ['quantity', 'name', 'short_description', 'replacement_value', 'resale_value', 'total_value',
    'manufacturer', 'model', 'category', 'room', 'condition', 'purchase_source', 'purchase_price', 'purchase_date', 'serial_number', 'website_url', 'manual_url', 'info_complete', 'details'];
const numericFields = new Set(['quantity', 'replacement_value', 'resale_value', 'purchase_price']);
const moneyFields = new Set(['replacement_value', 'resale_value', 'total_value']);

const filterSelects = new Set([ui.categoryFilter, ui.photoFilter, ui.infoFilter, ui.roomFilter]);
let numFilters = 0;

let items = [];
let selectedItem = null;
let selectedCategory = '';
let selectedPhotoFilter = '';
let selectedRoom = '';
let selectedInfoFilter = '';
const selectedIDs = new Set();
let searchText = "";
let saving = false;
let deleting = false;

function fieldLabel(key) {
    if (key === 'quantity') return 'Qty';
    if (key === 'number') return 'Num';
    if (key === 'room') return 'Room / Location';
    if (key === 'info_complete') return 'Information Complete';
    if (key === 'website_url') return 'Website URL';
    if (key === 'manual_url') return 'Manual URL';
    if (key === 'short_description') return 'Short Description';
    if (key === 'replacement_value') return 'Replacement Value';
    if (key === 'resale_value') return 'Resale Value';
    return key.replaceAll('_', ' ').replace(/\b\w/g, char => char.toUpperCase());
}

function money(amount) {
    return '$' + Number(amount || 0).toLocaleString('en-US');
}

function updateFilterState() {
    let count = 0;
    filterSelects.forEach((el) => {
        const active = el.value !== "";
        el.classList.toggle('is-active', active);
        if (active) count++;
    });
    ui.filterButton.classList.toggle('is-active', count !== 0);
    ui.clearFilters.disabled = count === 0;
    ui.filtersCount.classList.toggle('is-hidden', count === 0);
    ui.filtersCount.innerText = count === 1 ? String(count) + " filter set" : String(count) + " filters set";
    ui.filtersCount.classList.toggle('is-hidden', count === 0);
}

function clearFilters() {
    // reset all the select elements to ''
    filterSelects.forEach((el) => {
        el.value = "";
    });
    selectedCategory = selectedPhotoFilter = selectedInfoFilter = selectedRoom = '';
    renderTable();
}

function visibleItems() {
    updateFilterState();
    return items.filter(item => {
        const matchesCategory =
            !selectedCategory || item.category === selectedCategory;
        const matchesRoom = !selectedRoom || item.room === selectedRoom;
        const matchesInfo = selectedInfoFilter === '' || (selectedInfoFilter === 'complete' ? item.info_complete : !item.info_complete);

        const matchesSearch = [
            item.name,
            item.short_description,
            item.manufacturer,
            item.model,
            item.serial_number,
            item.details,
            item.room,
            item.purchase_source,
            item.website_url,
            item.manual_url
        ].some(value =>
            String(value ?? "").toLowerCase().includes(searchText)
        );

        const photoCount = Number(item.photo_count || 0);
        const matchesPhotos = selectedPhotoFilter === '' ||
            (selectedPhotoFilter === 'without' && photoCount === 0) ||
            (selectedPhotoFilter === 'with' && photoCount > 0);
        return matchesCategory && matchesRoom && matchesSearch && matchesPhotos && matchesInfo;
    }).sort((a, b) =>
        a.name.localeCompare(b.name) ||
        a.short_description.localeCompare(b.short_description)
    );
}

function renderTable() {
    const displayed = visibleItems();
    const header = document.createElement('tr');
    const selectHead = document.createElement('th');
    const selectAll = document.createElement('input');
    selectAll.type = 'checkbox';
    selectAll.setAttribute('aria-label', 'Select all visible items');
    selectAll.checked = displayed.length > 0 && displayed.every(item => selectedIDs.has(item.id));
    selectAll.indeterminate = displayed.some(item => selectedIDs.has(item.id)) && !selectAll.checked;
    selectAll.addEventListener('change', () => {
        displayed.forEach(item => selectAll.checked ? selectedIDs.add(item.id) : selectedIDs.delete(item.id));
        renderTable();
    });
    selectHead.append(selectAll);
    header.append(selectHead);
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
        const selectCell = document.createElement('td');
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.checked = selectedIDs.has(item.id);
        checkbox.setAttribute('aria-label', `Select ${item.name}`);
        checkbox.addEventListener('click', event => event.stopPropagation());
        checkbox.addEventListener('change', () => {
            checkbox.checked ? selectedIDs.add(item.id) : selectedIDs.delete(item.id);
            renderTable();
        });
        selectCell.append(checkbox);
        row.append(selectCell);
        for (const key of tableFields) {
            const td = document.createElement('td');
            if (moneyFields.has(key)) {
                td.classList.add('inventory-table__money');
                td.textContent = money(item[key]);
            } else {
                td.textContent = key === 'info_complete' ? (item.info_complete ? 'Yes' : 'Needs info') : String(key === 'number' ? index + 1 : (item[key] ?? ''));
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

    ui.bulk.textContent = `Bulk Edit (${selectedIDs.size})`;
    ui.bulk.disabled = selectedIDs.size === 0;
    ui.count.textContent = displayed.length;
    ui.value.textContent = money(displayed.reduce((sum, item) => sum + Number(item.total_value || 0), 0));
}

function createFields() {
    const category = ui.categoryFilter.cloneNode(true);
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
        if (key === "info_complete") {
            input = document.createElement("input");
            input.type = "checkbox";
        } else if (key === "category") {
            input = category;
        } else if (key === "details") {
            input = document.createElement("textarea");
            input.rows = 3;
        } else {
            input = document.createElement("input");
            input.type = numericFields.has(key) ? 'number' : 'text';
        }
        if (key === 'room') input.setAttribute('list', 'room-options');
        if (key === 'purchase_source') input.setAttribute('list', 'source-options');
        if (key === 'condition') {
            const choices = document.createElement('datalist');
            choices.id = 'condition-options';
            ['New', 'Excellent', 'Good', 'Fair', 'Poor'].forEach(v => {
                const option = document.createElement('option');
                option.value = v;
                choices.append(option);
            });
            if (!document.getElementById(choices.id)) document.body.append(choices);
            input.setAttribute('list', choices.id);
        }
        if (key === 'website_url' || key === 'manual_url') {
            input.type = 'url';
            input.placeholder = 'https://';
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
        if (key === "website_url" || key === "manual_url") {
            const link = document.createElement("a");
            link.className = "item-fields__link is-hidden";
            link.target = "_blank";
            link.rel = "noopener noreferrer";
            link.textContent = "Open link";
            link.dataset.forField = key;
            field.append(link);
        }
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
    updateItemLinks();
    ui.photoUploadLabel.classList.toggle('is-hidden', !editing || adding || !selectedItem?.id);
    ui.photoGallery.querySelectorAll('.photo-delete').forEach(button => button.classList.toggle('is-hidden', !editing));
    if (adding) {
        ui.dialogTitle.textContent = "Add Item";
    } else if (editing) {
        ui.dialogTitle.textContent = "Edit Item";
    } else {
        ui.dialogTitle.textContent = "Item Details";
    }
}

function openItem(item) {
    selectedItem = item;
    showError();
    for (const key of detailFields) {
        if (key === 'info_complete') ui.form.elements.namedItem(key).checked = Boolean(item[key]);
        else ui.form.elements.namedItem(key).value = item[key] ?? '';
    }
    setMode(false);
    updateItemLinks();
    loadPhotos(item.id);
    ui.dialog.classList.remove('is-hidden');
    ui.edit.focus({preventScroll: true});
    ui.dialog.querySelector('.item-dialog__panel').scrollTop = 0;
}

function closeItem() {
    if (saving) return;
    ui.dialog.classList.add('is-hidden');
    selectedItem = null;
}

function addItem(event) {
    selectedItem = {id: 0};
    showError();
    ui.form.reset();
    for (const key of detailFields) {
        if (key === 'info_complete') ui.form.elements.namedItem(key).checked = false;
        else ui.form.elements.namedItem(key).value = key === 'purchase_price' ? '0' : '';
    }
    ui.photoGallery.replaceChildren();
    ui.photoStatus.textContent = 'Save the item before adding photos.';
    setMode(true, true);
    updateItemLinks();
    ui.dialog.classList.remove('is-hidden');
    ui.edit.focus();
    ui.form.elements.namedItem('name').focus();
}

function readItem() {
    const item = {id: selectedItem?.id ?? 0};
    for (const key of detailFields) {
        if (key === 'total_value') continue; // Computed by the server.
        if (key === 'info_complete') {
            item[key] = ui.form.elements.namedItem(key).checked;
            continue;
        }
        const value = ui.form.elements.namedItem(key).value.trim();
        if (numericFields.has(key)) {
            if (!(key === 'purchase_price' && value === '') && !/^\d+$/.test(value)) throw new Error(`${fieldLabel(key)} must be a nonnegative whole number.`);
            item[key] = Number(value || 0);
            if (!Number.isSafeInteger(item[key])) throw new Error(`${fieldLabel(key)} is too large.`);
        } else if (key === 'website_url' || key === 'manual_url') {
            if (value) {
                let parsed;
                try {
                    parsed = new URL(value);
                } catch {
                    throw new Error(fieldLabel(key) + ' must be a valid URL.');
                }
                if (!['http:', 'https:'].includes(parsed.protocol)) throw new Error(fieldLabel(key) + ' must use http or https.');
            }
            item[key] = value;
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
            headers: {'Content-Type': 'application/json'},
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
    const payload = {id: selectedItem.id};
    deleting = true;
    ui.save.disabled = true;
    ui.delete.disabled = true;
    ui.confirmDelete.disabled = true;
    showError();
    try {
        const response = await fetch('/api/delete', {
            method: 'DELETE',
            headers: {'Content-Type': 'application/json'},
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
    searchText = event.target.value.trim().toLowerCase();
    renderTable();
}

async function loadItems() {
    const response = await fetch('/api/items');
    if (!response.ok) throw new Error(`Unable to load items (HTTP ${response.status}).`);
    items = await response.json();
    const knownIDs = new Set(items.map(item => item.id));
    for (const id of selectedIDs) if (!knownIDs.has(id)) selectedIDs.delete(id);
    const roomOptions = [...new Set(items.map(item => String(item.room || '').trim()).filter(Boolean))].sort();
    ui.roomFilter.replaceChildren(new Option('All', ''), ...roomOptions.map(room => new Option(room, room)));
    ui.roomFilter.value = selectedRoom;
    const list = document.getElementById('room-options');
    list.replaceChildren(...roomOptions.map(room => new Option(room, room)));
    renderTable();
}


function updateItemLinks() {
    for (const key of ['website_url', 'manual_url']) {
        const anchor = ui.fields.querySelector('[data-for-field="' + key + '"]');
        if (!anchor) continue;
        const value = ui.form.elements.namedItem(key)?.value.trim() || '';
        let valid = false;
        try {
            const parsed = new URL(value);
            valid = ['https:', 'http:'].includes(parsed.protocol);
        } catch {
        }
        anchor.classList.toggle('is-hidden', !valid);
        if (valid) anchor.href = value;
        else anchor.removeAttribute('href');
    }
}

function openBulkEditor() {
    if (!selectedIDs.size) return;
    ui.bulkCount.textContent = 'Update ' + selectedIDs.size + ' selected item(s). Only the chosen field will change.';
    ui.bulkField.value = 'room';
    ui.bulkValue.value = '';
    ui.bulkValue.type = 'text';
    ui.bulkValue.placeholder = '';
    ui.bulkError.textContent = '';
    ui.bulkValue.setAttribute('list', 'room-options');
    ui.bulkDialog.showModal();
    ui.bulkValue.focus();
}

async function applyBulkEdit() {
    const ids = [...selectedIDs];
    const field = ui.bulkField.value;
    const value = ui.bulkField.value === 'info_complete' ? ui.bulkValue.value : ui.bulkValue.value.trim();
    if (!ids.length) return;
    if (!confirm('Set ' + fieldLabel(field) + ' to ' + (value || '(blank)') + ' on ' + ids.length + ' items?')) return;
    ui.bulkApply.disabled = true;
    ui.bulkError.textContent = '';
    try {
        const response = await fetch('/api/items/bulk', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ids, field, value})
        });
        if (!response.ok) throw new Error((await response.text()).trim() || 'Bulk update failed.');
        selectedIDs.clear();
        await loadItems();
        ui.bulkDialog.close();
    } catch (error) {
        ui.bulkError.textContent = error.message;
        await loadItems().catch(() => {
        });
    } finally {
        ui.bulkApply.disabled = false;
    }
}

function init() {
    createFields();
    ui.filterButton.addEventListener('click', () => {
        ui.toolbarFilters.classList.toggle('is-hidden');
    });
    ui.categoryFilter.addEventListener('change', () => {
        selectedCategory = ui.categoryFilter.value;
        renderTable();
    });
    ui.photoFilter.addEventListener('change', () => {
        selectedPhotoFilter = ui.photoFilter.value;
        renderTable();
    });
    ui.infoFilter.addEventListener('change', () => {
        selectedInfoFilter = ui.infoFilter.value;
        renderTable();
    });
    ui.roomFilter.addEventListener('change', () => {
        selectedRoom = ui.roomFilter.value;
        renderTable();
    });
    ui.clearFilters.addEventListener('click', clearFilters)
    ui.bulk.addEventListener('click', openBulkEditor);
    ui.bulkCancel.addEventListener('click', () => ui.bulkDialog.close());
    ui.bulkApply.addEventListener('click', applyBulkEdit);
    ui.bulkField.addEventListener('change', () => {
        const field = ui.bulkField.value;
        ui.bulkValue.value = '';
        if (field === 'room') ui.bulkValue.setAttribute('list', 'room-options');
        else if (field === 'purchase_source') ui.bulkValue.setAttribute('list', 'source-options');
        else if (field === 'condition') ui.bulkValue.setAttribute('list', 'condition-options');
        else ui.bulkValue.removeAttribute('list');
        if (field === 'info_complete') {
            ui.bulkValue.type = 'text';
            ui.bulkValue.placeholder = 'true (complete) or false (needs info)';
            ui.bulkValue.setAttribute('list', 'completion-options');
        } else {
            ui.bulkValue.type = field.endsWith('_url') ? 'url' : 'text';
            ui.bulkValue.placeholder = '';
        }
    });
    ui.fields.addEventListener('input', event => {
        if (['website_url', 'manual_url'].includes(event.target.name)) updateItemLinks();
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
                    const result = await fetch(`/api/items/${itemID}/photos/${photo.id}`, {method: 'DELETE'});
                    if (!result.ok) throw new Error('Unable to delete photo.');
                    await loadPhotos(itemID);
                    await loadItems();
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
        ui.photoInput.value = '';
        return;
    }
    const form = new FormData();
    form.append('photo', file);
    ui.photoInput.disabled = true;
    ui.photoStatus.textContent = 'Uploading photo…';
    try {
        const response = await fetch(`/api/items/${itemID}/photos`, {method: 'POST', body: form});
        if (!response.ok) throw new Error(await response.text() || 'Upload failed.');
        await loadPhotos(itemID);
        await loadItems();
    } catch (error) {
        ui.photoStatus.textContent = error.message;
    } finally {
        ui.photoInput.value = '';
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
