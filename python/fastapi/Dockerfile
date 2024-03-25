FROM python:3.10-slim

RUN mkdir /app
EXPOSE 8000

WORKDIR /app
COPY requirements.txt /app
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY main.py /app

CMD ["uvicorn", "--host", "0.0.0.0", "--port", "8000", "main:app"]
