FROM python:alpine

WORKDIR /app
COPY load_generator.py .
RUN pip3 install requests

CMD ["python", "load_generator.py"]