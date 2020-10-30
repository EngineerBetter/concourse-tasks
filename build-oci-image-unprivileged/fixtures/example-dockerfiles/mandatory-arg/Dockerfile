FROM alpine
ARG MANDATORY_ARG
ARG OTHER_ARG

RUN \
    # Check for mandatory build arguments
        : "${MANDATORY_ARG:?Build argument needs to be set and non-empty.}"

RUN \
    # Check for mandatory build arguments
        : "${OTHER_ARG:?Build argument needs to be set and non-empty.}"