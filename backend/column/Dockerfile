FROM python:3.12-slim AS lint

WORKDIR /app
COPY . /app

# black をインストール
RUN pip install --no-cache-dir black
RUN black . --check --exclude="(venv|proto|gen)"

# flake8 をインストール
RUN pip install --no-cache-dir flake8
RUN flake8 . --exclude=venv,proto,gen --ignore=E501

FROM python:3.12-slim AS runtime

WORKDIR /app
COPY . /app

# 本番用の依存関係をインストール
RUN pip install --no-cache-dir -r requirements.txt

ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app:/app/gen

CMD ["python", "./cmd/server/main.py", "server"]
