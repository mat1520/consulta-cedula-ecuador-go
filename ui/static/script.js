// Elementos del DOM
const form = document.getElementById('cedula-form');
const cedulaInput = document.getElementById('cedula-input');
const submitBtn = document.getElementById('submit-btn');
const spinner = document.getElementById('spinner');
const btnText = document.querySelector('.btn-text');
const resultsContainer = document.getElementById('results-container');

// Estado de la aplicación
let isLoading = false;

// Event listener para el formulario
form.addEventListener('submit', async (event) => {
    event.preventDefault();
    
    // Evitar múltiples envíos
    if (isLoading) return;
    
    const cedula = cedulaInput.value.trim();
    
    // Validación frontend
    if (!validarCedulaFrontend(cedula)) {
        mostrarError('Por favor, ingresa un número de cédula válido de 10 dígitos.');
        return;
    }
    
    // Realizar la consulta
    await consultarCedula(cedula);
});

// Validación en el frontend
function validarCedulaFrontend(cedula) {
    // Verificar que tenga exactamente 10 dígitos
    const regexCedula = /^[0-9]{10}$/;
    return regexCedula.test(cedula);
}

// Función principal para consultar la cédula
async function consultarCedula(cedula) {
    try {
        // Mostrar estado de carga
        mostrarCarga(true);
        
        // Realizar petición al backend
        const response = await fetch('/api/consultar', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ cedula: cedula })
        });
        
        // Parsear respuesta JSON
        const data = await response.json();
        
        // Manejar respuesta según el código de estado
        if (response.ok) {
            // Consulta exitosa (200)
            mostrarExito(data);
        } else if (response.status === 404) {
            // Cédula no encontrada (404)
            mostrarError('❌ Cédula no encontrada en los registros públicos.');
        } else if (response.status === 400) {
            // Error de validación (400)
            mostrarError(`⚠️ Error de validación: ${data.error}`);
        } else if (response.status === 500) {
            // Error interno del servidor (500)
            mostrarError('🔧 Error interno del servidor. Inténtalo nuevamente más tarde.');
        } else {
            // Otros errores
            mostrarError(`❌ Error inesperado: ${data.error || 'Error desconocido'}`);
        }
        
    } catch (error) {
        // Error de red o conexión
        console.error('Error en la petición:', error);
        mostrarError('🌐 Error de conexión. Verifica tu conexión a internet e inténtalo nuevamente.');
        
    } finally {
        // Ocultar estado de carga
        mostrarCarga(false);
    }
}

// Mostrar/ocultar estado de carga
function mostrarCarga(mostrar) {
    isLoading = mostrar;
    
    if (mostrar) {
        submitBtn.disabled = true;
        submitBtn.classList.add('loading');
        btnText.style.opacity = '0';
        spinner.style.display = 'block';
    } else {
        submitBtn.disabled = false;
        submitBtn.classList.remove('loading');
        btnText.style.opacity = '1';
        spinner.style.display = 'none';
    }
}

// Mostrar resultado exitoso
function mostrarExito(data) {
    resultsContainer.innerHTML = `
        <div class="result-success">
            <div class="result-title">✅ ¡Consulta Exitosa!</div>
            <div class="result-data">
                <p><strong>Nombre:</strong> ${data.nombre}</p>
                <p><strong>Apellido:</strong> ${data.apellido}</p>
                <p><strong>Cédula:</strong> ${cedulaInput.value}</p>
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

// Validación en tiempo real del input
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

// Permitir envío con Enter
cedulaInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
        event.preventDefault();
        form.dispatchEvent(new Event('submit'));
    }
});

// Enfocar automáticamente el input al cargar la página
document.addEventListener('DOMContentLoaded', () => {
    cedulaInput.focus();
});

// Función de utilidad para limpiar y reiniciar el formulario
function limpiarFormulario() {
    cedulaInput.value = '';
    resultsContainer.style.display = 'none';
    cedulaInput.focus();
}

// Agregar listener para doble clic en el input (limpiar)
cedulaInput.addEventListener('dblclick', limpiarFormulario);
