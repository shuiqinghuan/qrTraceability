FROM python:3.11-slim

LABEL maintainer="Agricultural Traceability System"
LABEL description="Django backend for Agricultural Traceability System"

WORKDIR /app

RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY backend/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY backend/ .
COPY backend/manage.py .

RUN mkdir -p /app/logs /app/staticfiles /app/media

RUN python manage.py collectstatic --noinput || true

EXPOSE 8000

ENV PYTHONUNBUFFERED=1
ENV DJANGO_SETTINGS_MODULE=backend.settings

CMD ["gunicorn", "backend.wsgi:application", "--bind", "0.0.0.0:8000", "--workers", "3"]
