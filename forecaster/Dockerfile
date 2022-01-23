FROM python:3.9

RUN pip3 install requests flask pmdarima pandas

WORKDIR /app

COPY forecaster.py .

EXPOSE 5000

CMD ["python", "forecaster.py"]