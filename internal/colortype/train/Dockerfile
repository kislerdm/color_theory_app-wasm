FROM python:3.8

WORKDIR /app

COPY main.py main.py
COPY requirements.txt requirements.txt

RUN pip install --upgrade pip \
    &&  pip install -r requirements.txt

ENTRYPOINT ["python", "main.py"]
