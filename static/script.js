// Глобальная переменная для хранения ID опроса, который нужно удалить
let surveyToDelete = null;

function addQuestion() {
    const template = document.getElementById('questionTemplate');
    const container = document.getElementById('questionContainer');
    const clone = template.content.cloneNode(true);
    const newQuestion = clone.querySelector('.question-block');
    
    // Добавляем с анимацией
    newQuestion.style.opacity = '0';
    container.appendChild(clone);
    
    // Запускаем анимацию появления
    requestAnimationFrame(() => {
        newQuestion.style.opacity = '1';
        newQuestion.style.transform = 'translateY(0)';
    });
    
    updateQuestionNumbers();
}

function handleQuestionTypeChange(select) {
    const optionsContainer = select.closest('.card-body').querySelector('.options-container');
    if (select.value === 'choice' || select.value === 'multiple') {
        optionsContainer.style.display = 'block';
        optionsContainer.style.animation = 'fadeIn 0.3s ease forwards';
    } else {
        optionsContainer.style.animation = 'fadeOut 0.3s ease forwards';
        setTimeout(() => {
            optionsContainer.style.display = 'none';
        }, 300);
    }
}

function addOption(button) {
    const optionsContainer = button.closest('.options-container');
    const optionCount = optionsContainer.querySelectorAll('.input-group').length;
    const newOption = document.createElement('div');
    newOption.className = 'mb-2';
    newOption.innerHTML = `
        <div class="input-group">
            <input type="text" class="form-control" placeholder="Вариант ${optionCount + 1}">
            <button type="button" class="btn btn-outline-danger hover-scale" onclick="removeOption(this)">
                <i class="fas fa-times"></i>
            </button>
        </div>
    `;
    
    // Добавляем с анимацией
    newOption.style.opacity = '0';
    optionsContainer.insertBefore(newOption, button);
    
    requestAnimationFrame(() => {
        newOption.style.opacity = '1';
        newOption.style.animation = 'slideDown 0.3s ease forwards';
    });
}

function removeOption(button) {
    const optionDiv = button.closest('.mb-2');
    optionDiv.style.animation = 'fadeOut 0.3s ease forwards';
    setTimeout(() => {
        optionDiv.remove();
    }, 300);
}

function updateQuestionNumbers() {
    const questions = document.querySelectorAll('.question-block h5');
    questions.forEach((question, index) => {
        if (!question.closest('.template')) {
            question.textContent = `Вопрос ${index + 1}`;
        }
    });
}

function deleteQuestion(questionId) {
    const questionBlock = typeof questionId === 'object' 
        ? questionId.closest('.question-block')
        : document.querySelector(`.question-block[data-question-id="${questionId}"]`);

    if (!confirm('Вы уверены, что хотите удалить этот вопрос?')) {
        return;
    }

    // Добавляем анимацию удаления
    questionBlock.style.animation = 'fadeOut 0.3s ease forwards';
    
    setTimeout(() => {
        if (typeof questionId === 'object') {
            questionBlock.remove();
            updateQuestionNumbers();
            return;
        }

        fetch(`/api/question/${questionId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                questionBlock.remove();
                updateQuestionNumbers();
            } else {
                alert('Ошибка при удалении вопроса: ' + data.error);
                // Возвращаем блок в исходное состояние
                questionBlock.style.animation = '';
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Ошибка при удалении вопроса. Пожалуйста, попробуйте снова.');
            // Возвращаем блок в исходное состояние
            questionBlock.style.animation = '';
        });
    }, 300);
}

function deleteSurvey(surveyId) {
    surveyToDelete = surveyId;
    const modal = new bootstrap.Modal(document.getElementById('deleteConfirmModal'));
    modal.show();
}

function saveSurvey() {
    const questions = [];
    document.querySelectorAll('.question-block').forEach(block => {
        if (block.matches('.template')) return;

        const questionData = {
            content: block.querySelector('input[type="text"]').value.trim(),
            type: block.querySelector('select').value,
            options: []
        };

        if (!questionData.content) {
            alert('Пожалуйста, заполните текст всех вопросов');
            return;
        }

        if (questionData.type === 'choice' || questionData.type === 'multiple') {
            const options = [];
            block.querySelectorAll('.options-container .input-group input[type="text"]').forEach(input => {
                const value = input.value.trim();
                if (value) {
                    options.push(value);
                }
            });

            if (options.length < 2) {
                alert('Для вопросов с выбором необходимо указать как минимум 2 варианта ответа');
                return;
            }

            questionData.options = options;
        }

        questions.push(questionData);
    });

    if (questions.length === 0) {
        alert('Добавьте хотя бы один вопрос');
        return;
    }

    // Добавляем анимацию загрузки
    const saveButton = document.querySelector('button[onclick="saveSurvey()"]');
    const originalText = saveButton.innerHTML;
    saveButton.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Сохранение...';
    saveButton.disabled = true;

    fetch(window.location.href, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ questions: questions })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            window.location.href = '/';
        } else {
            alert('Ошибка при сохранении опроса: ' + data.error);
            saveButton.innerHTML = originalText;
            saveButton.disabled = false;
        }
    })
    .catch(error => {
        console.error('Ошибка:', error);
        alert('Ошибка при сохранении опроса. Пожалуйста, попробуйте снова.');
        saveButton.innerHTML = originalText;
        saveButton.disabled = false;
    });
}

// Обработчик подтверждения удаления
document.addEventListener('DOMContentLoaded', function() {
    const confirmDeleteBtn = document.getElementById('confirmDelete');
    if (confirmDeleteBtn) {
        confirmDeleteBtn.addEventListener('click', function() {
            if (surveyToDelete) {
                const surveyCard = document.querySelector(`.survey-card[data-survey-id="${surveyToDelete}"]`);
                
                fetch(`/api/survey/${surveyToDelete}`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        if (surveyCard) {
                            const cardContainer = surveyCard.closest('.col-md-4');
                            cardContainer.style.animation = 'fadeOut 0.3s ease forwards';
                            setTimeout(() => {
                                cardContainer.remove();
                                // Проверяем, остались ли опросы
                                const remainingSurveys = document.querySelectorAll('.survey-card');
                                if (remainingSurveys.length === 0) {
                                    const row = document.querySelector('.row');
                                    if (row) {
                                        const noSurveysDiv = document.createElement('div');
                                        noSurveysDiv.className = 'col-12 fade-in';
                                        noSurveysDiv.innerHTML = '<p>Нет доступных опросов.</p>';
                                        row.appendChild(noSurveysDiv);
                                    }
                                }
                            }, 300);
                        }
                        const modal = bootstrap.Modal.getInstance(document.getElementById('deleteConfirmModal'));
                        modal.hide();
                    } else {
                        alert('Ошибка при удалении опроса: ' + data.error);
                    }
                })
                .catch(error => {
                    console.error('Ошибка:', error);
                    alert('Ошибка при удалении опроса. Пожалуйста, попробуйте снова.');
                });
            }
        });
    }
});
