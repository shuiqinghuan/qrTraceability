FROM python:3.11-slim

WORKDIR /app

# 复制依赖文件
COPY backend/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 复制项目文件
COPY manage.py .
COPY backend/ ./backend/
COPY products/ ./products/

# 创建必要的目录
RUN mkdir -p staticfiles media logs

# 设置Python路径
ENV PYTHONPATH=/app

EXPOSE 8000

CMD ["gunicorn", "backend.wsgi:application", "--bind", "0.0.0.0:8000", "--workers", "2"]
