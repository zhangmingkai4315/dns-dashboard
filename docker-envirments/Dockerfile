FROM  sameersbn/bind:9.9.5-20170626
ADD ./root-conf/* /etc/bind/
EXPOSE 53/udp
RUN mkdir /var/cache/bind/dnslog
RUN chown bind.bind -R /var/cache/bind/dnslog
CMD ["bash","-c","named ; tail -f /dev/null" ]