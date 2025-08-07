// Elementos del DOM
const tabButtons = document.querySelectorAll('.tab-button');
const cedulaForm = document.getElementById('cedula-form');
const nombresForm = document.getElementById('nombres-form');
const cedulaInput = document.getElementById('cedula-input');
const nombresInput = document.getElementById('nombres-input');
const apellidosInput = document.getElementById('apellidos-input');
const submitCedulaBtn = document.getElementById('submit-cedula-btn');
const submitNombresBtn = document.getElementById('submit-nombres-btn');
const resultsContainer = document.getElementById('results-container');

// Estado de la aplicación
let isLoading = false;

// Event listeners para las pestañas
tabButtons.forEach(button => {
    button.addEventListener('click', () => {
        const tab = button.dataset.tab;
        switchTab(tab);
    });
});

// Event listener para el formulario de cédula
cedulaForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    
    if (isLoading) return;
    
    const cedula = cedulaInput.value.trim();
    
    if (!validarCedulaFrontend(cedula)) {
        mostrarError('Por favor, ingresa un número de cédula válido de 10 dígitos.');
        return;
    }
    
    await consultarCedula(cedula);
});

// Event listener para el formulario de nombres
nombresForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    
    if (isLoading) return;
    
    const nombres = nombresInput.value.trim();
    const apellidos = apellidosInput.value.trim();
    
    if (!validarNombresFrontend(nombres, apellidos)) {
        mostrarError('Por favor, ingresa nombres y apellidos válidos.');
        return;
    }
    
    await consultarPorNombres(nombres, apellidos);
});

// Función para cambiar entre pestañas
function switchTab(tab) {
    // Actualizar botones de pestaña
    tabButtons.forEach(btn => {
        btn.classList.toggle('active', btn.dataset.tab === tab);
    });
    
    // Mostrar/ocultar formularios
    cedulaForm.classList.toggle('active', tab === 'cedula');
    nombresForm.classList.toggle('active', tab === 'nombres');
    
    // Limpiar resultados
    resultsContainer.style.display = 'none';
    
    // Enfocar el input correspondiente
    if (tab === 'cedula') {
        cedulaInput.focus();
    } else {
        nombresInput.focus();
    }
}

// Validación en el frontend para cédulas
function validarCedulaFrontend(cedula) {
    const regexCedula = /^[0-9]{10}$/;
    return regexCedula.test(cedula);
}

// Validación en el frontend para nombres
function validarNombresFrontend(nombres, apellidos) {
    const regexNombres = /^[A-Za-zÀ-ÿ\u00f1\u00d1\s]+$/;
    return nombres.length >= 2 && apellidos.length >= 2 && 
           regexNombres.test(nombres) && regexNombres.test(apellidos);
}

// Función principal para consultar la cédula
async function consultarCedula(cedula) {
    try {
        mostrarCarga(true, 'cedula');
        
        const response = await fetch('/api/consultar', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ cedula: cedula })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            mostrarExitoCedula(data, cedula);
        } else if (response.status === 404) {
            mostrarError('❌ Cédula no encontrada en los registros públicos.');
        } else if (response.status === 400) {
            mostrarError(`⚠️ Error de validación: ${data.error}`);
        } else if (response.status === 500) {
            mostrarError('🔧 Error interno del servidor. Inténtalo nuevamente más tarde.');
        } else {
            mostrarError(`❌ Error inesperado: ${data.error || 'Error desconocido'}`);
        }
        
    } catch (error) {
        console.error('Error en la petición:', error);
        mostrarError('🌐 Error de conexión. Verifica tu conexión a internet e inténtalo nuevamente.');
        
    } finally {
        mostrarCarga(false, 'cedula');
    }
}

// Función principal para consultar por nombres
async function consultarPorNombres(nombres, apellidos) {
    try {
        mostrarCarga(true, 'nombres');
        
        const response = await fetch('/api/consultar-nombres', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ nombres: nombres, apellidos: apellidos })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            mostrarExitoNombres(data);
        } else if (response.status === 404) {
            mostrarAlternativasLegales(nombres, apellidos);
        } else if (response.status === 400) {
            mostrarError(`⚠️ Error de validación: ${data.error}`);
        } else if (response.status === 500) {
            // El error del servidor contiene información sobre alternativas legales
            mostrarAlternativasLegales(nombres, apellidos, data.error);
        } else {
            mostrarError(`❌ Error inesperado: ${data.error || 'Error desconocido'}`);
        }
        
    } catch (error) {
        console.error('Error en la petición:', error);
        mostrarError('🌐 Error de conexión. Verifica tu conexión a internet e inténtalo nuevamente.');
        
    } finally {
        mostrarCarga(false, 'nombres');
    }
}

// Mostrar/ocultar estado de carga
function mostrarCarga(mostrar, tipo = 'cedula') {
    isLoading = mostrar;
    
    const btn = tipo === 'cedula' ? submitCedulaBtn : submitNombresBtn;
    const btnText = btn.querySelector('.btn-text');
    const spinner = btn.querySelector('.spinner');
    
    if (mostrar) {
        btn.disabled = true;
        btn.classList.add('loading');
        btnText.style.opacity = '0';
        spinner.style.display = 'block';
    } else {
        btn.disabled = false;
        btn.classList.remove('loading');
        btnText.style.opacity = '1';
        spinner.style.display = 'none';
    }
}

// Mostrar resultado exitoso para búsqueda por cédula
function mostrarExitoCedula(data, cedula) {
    resultsContainer.innerHTML = `
        <div class="result-success">
            <div class="result-title">✅ ¡Consulta Exitosa!</div>
            <div class="result-data">
                <p><strong>Cédula:</strong> ${cedula}</p>
                <p><strong>Nombre:</strong> ${data.nombre}</p>
                <p><strong>Apellido:</strong> ${data.apellido}</p>
            </div>
        </div>
    `;
    
    resultsContainer.style.display = 'block';
    scrollToResults();
}

// Mostrar resultado exitoso para búsqueda por nombres
function mostrarExitoNombres(data) {
    resultsContainer.innerHTML = `
        <div class="result-success">
            <div class="result-title">✅ ¡Cédula Encontrada!</div>
            <div class="result-data">
                <p><strong>Nombres:</strong> ${data.nombres}</p>
                <p><strong>Apellidos:</strong> ${data.apellidos}</p>
                <p><strong>Cédula:</strong> ${data.cedula}</p>
            </div>
        </div>
    `;
    
    resultsContainer.style.display = 'block';
    scrollToResults();
}

// Mostrar mensaje de error
function mostrarError(mensaje) {
    resultsContainer.innerHTML = `
        <div class="result-error">
            <div class="result-title">Error en la Consulta</div>
            <div class="result-data">
                <p>${mensaje}</p>
            </div>
        </div>
    `;
    
    resultsContainer.style.display = 'block';
    scrollToResults();
}

// Hacer scroll suave hacia los resultados
function scrollToResults() {
    setTimeout(() => {
        resultsContainer.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        });
    }, 100);
}

// Validación en tiempo real del input de cédula
cedulaInput.addEventListener('input', (event) => {
    let value = event.target.value;
    
    // Permitir solo números
    value = value.replace(/[^0-9]/g, '');
    
    // Limitar a 10 caracteres
    if (value.length > 10) {
        value = value.slice(0, 10);
    }
    
    event.target.value = value;
    
    // Ocultar resultados anteriores si se modifica la cédula
    if (resultsContainer.style.display === 'block') {
        resultsContainer.style.display = 'none';
    }
});

// Validación en tiempo real del input de nombres
nombresInput.addEventListener('input', (event) => {
    let value = event.target.value;
    
    // Permitir solo letras, espacios y acentos
    value = value.replace(/[^A-Za-zÀ-ÿ\u00f1\u00d1\s]/g, '');
    
    event.target.value = value;
    
    // Ocultar resultados anteriores
    if (resultsContainer.style.display === 'block') {
        resultsContainer.style.display = 'none';
    }
});

// Validación en tiempo real del input de apellidos
apellidosInput.addEventListener('input', (event) => {
    let value = event.target.value;
    
    // Permitir solo letras, espacios y acentos
    value = value.replace(/[^A-Za-zÀ-ÿ\u00f1\u00d1\s]/g, '');
    
    event.target.value = value;
    
    // Ocultar resultados anteriores
    if (resultsContainer.style.display === 'block') {
        resultsContainer.style.display = 'none';
    }
});

// Permitir envío con Enter en el input de cédula
cedulaInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        cedulaForm.dispatchEvent(new Event('submit'));
    }
});

// Permitir envío con Enter en los inputs de nombres
nombresInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        apellidosInput.focus();
    }
});

apellidosInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        nombresForm.dispatchEvent(new Event('submit'));
    }
});

// Enfocar automáticamente el input al cargar la página
document.addEventListener('DOMContentLoaded', () => {
    cedulaInput.focus();
});

// Función de utilidad para limpiar y reiniciar formularios
function limpiarFormularios() {
    cedulaInput.value = '';
    nombresInput.value = '';
    apellidosInput.value = '';
    resultsContainer.style.display = 'none';
    
    // Enfocar según la pestaña activa
    const activeTab = document.querySelector('.tab-button.active').dataset.tab;
    if (activeTab === 'cedula') {
        cedulaInput.focus();
    } else {
        nombresInput.focus();
    }
}

// Agregar listener para doble clic en los inputs (limpiar)
cedulaInput.addEventListener('dblclick', limpiarFormularios);
nombresInput.addEventListener('dblclick', limpiarFormularios);
apellidosInput.addEventListener('dblclick', limpiarFormularios);

// Mostrar alternativas legales para consulta por nombres
function mostrarAlternativasLegales(nombres, apellidos, errorMessage = '') {
    resultsContainer.innerHTML = `
        <div class="result-alternatives">
            <div class="result-title">📋 Alternativas Legales Disponibles</div>
            <div class="search-info">
                <p><strong>Búsqueda solicitada:</strong> ${nombres} ${apellidos}</p>
                <p class="info-text">⚠️ La consulta por nombres no está disponible a través de APIs públicas gratuitas.</p>
            </div>
            
            <div class="alternatives-grid">
                <div class="alternative-card judicial">
                    <div class="card-header">
                        <span class="card-icon">🏛️</span>
                        <h3>Función Judicial (SATJE)</h3>
                    </div>
                    <div class="card-content">
                        <p>Consulta de procesos judiciales por nombre</p>
                        <ul>
                            <li>Buscar si una persona tiene procesos judiciales</li>
                            <li>Servicio oficial y gratuito</li>
                            <li>Búsqueda por nombres y apellidos</li>
                        </ul>
                        <a href="https://procesosjudiciales.funcionjudicial.gob.ec/busqueda" target="_blank" class="card-button">
                            Consultar SATJE
                        </a>
                    </div>
                </div>
                
                <div class="alternative-card electoral">
                    <div class="card-header">
                        <span class="card-icon">🗳️</span>
                        <h3>Consejo Nacional Electoral</h3>
                    </div>
                    <div class="card-content">
                        <p>Consulta de personas registradas para votar</p>
                        <ul>
                            <li>Búsqueda por nombre y apellido</li>
                            <li>Solo ciudadanos habilitados para elecciones</li>
                            <li>Servicio oficial gratuito</li>
                        </ul>
                        <a href="https://www.cne.gob.ec/" target="_blank" class="card-button">
                            Consultar CNE
                        </a>
                    </div>
                </div>
                
                <div class="alternative-card iess">
                    <div class="card-header">
                        <span class="card-icon">🏥</span>
                        <h3>IESS</h3>
                    </div>
                    <div class="card-content">
                        <p>Consulta de afiliados (protegida)</p>
                        <ul>
                            <li>Información de afiliación</li>
                            <li>Protegida con captcha</li>
                            <li>No tiene API pública</li>
                        </ul>
                        <a href="https://www.iess.gob.ec/" target="_blank" class="card-button">
                            Visitar IESS
                        </a>
                    </div>
                </div>
                
                <div class="alternative-card paid">
                    <div class="card-header">
                        <span class="card-icon">💰</span>
                        <h3>Servicios de Pago</h3>
                    </div>
                    <div class="card-content">
                        <p>EcuadorLegalOnline - Consulta completa</p>
                        <ul>
                            <li>Búsqueda por nombres y apellidos</li>
                            <li>Datos completos: cédula, estado civil, profesión</li>
                            <li>Servicio de pago con garantía</li>
                        </ul>
                        <a href="https://tramites.ecuadorlegalonline.com/" target="_blank" class="card-button paid-button">
                            Servicio de Pago
                        </a>
                    </div>
                </div>
            </div>
            
            <div class="recommendation">
                <div class="recommendation-header">
                    <span class="recommendation-icon">💡</span>
                    <h3>Recomendación</h3>
                </div>
                <p>Use el <strong>servicio de consulta por cédula</strong> que funciona con datos oficiales del SRI (gratuito y confiable).</p>
                <button onclick="switchTab('cedula')" class="switch-tab-button">
                    Ir a Consulta por Cédula
                </button>
            </div>
        </div>
    `;
    
    resultsContainer.style.display = 'block';
    scrollToResults();
}
